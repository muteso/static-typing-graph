package parser

import (
	"io"

	"fmt"
	"stg/template"
	"strconv"
	"sync"

	"gopkg.in/yaml.v2"
)

// Parses ALREADY opened template-file (or any another representation
// of it implementing io.Reader-interface) and returns result as
// TemplateHolder-struct; if any error occurs doesn't interrupt parsing
// and then returns error which contains all occured errors during parsing
//
// WARNING: template dont allow some graph design practices, as:
//   - maps nesting (which should be handled by making a new node/edge that
//     contains nested map etc.)
//   - arrays nesting (for the same reason as above)
//   - nesting maps within arrays and vise versa (for the same reason as above)
//   - map's and array's values (and keys in case of maps) can contain only the
//     same PRIMITIVE (int, string, etc) data types (for reasons of go data
//     types compatibility)
//   - labels can't nest within each other (which may lead to endless data-nesting
//     and implicit property- and connection-definitions which may cause hard-to-find
//     definition conflicts)
func ParseTemplate(file io.Reader) (*template.TemplateHolder, error) {
	b, err := io.ReadAll(file)
	if err != nil {
		return nil, parseError{
			"template",
			"can't read file: " + err.Error(),
		}
	}

	t := new(bTemplate)
	err = yaml.UnmarshalStrict(b, t)
	switch {
	case err != nil:
		return nil, parseError{
			"template",
			".yaml parsing error: " + err.Error(),
		}
	case len(t.BufEdges) == 0:
		return nil, parseError{
			"template",
			"there is no any edge definition",
		}
	case len(t.BufNodes) == 0:
		return nil, parseError{
			"template",
			"there is no any node definition",
		}
	}

	templ, err := t.toActual()
	if err != nil {
		return nil, err
	}
	return templ, nil
}

// Concurrently mutates bTemplate to actual Template-struct and return
// it; don't interrupts on error occurences but at end of the parsing
// process returns them all in single error-value
func (bt bTemplate) toActual() (*template.TemplateHolder, error) {
	c := newContext()

	done := new(sync.WaitGroup)
	// parses and transforms edges with their properties
	done.Add(len(bt.BufEdges))
	for k, v := range bt.BufEdges {
		go func(k string, v bEdge) {
			defer done.Done()
			v.nesting = append(v.nesting, "edges", k)
			v.toActual(c)
		}(k, v)
	}
	done.Wait()
	// parses and transforms labels with their properties
	done.Add(len(bt.BufLabels))
	for k, v := range bt.BufLabels {
		go func(k string, v bLabel) {
			defer done.Done()
			v.nesting = append(v.nesting, "labels", k)
			v.toActual(c)
		}(k, v)
	}
	done.Wait()
	// maps already parsed labels with each other using already parsed edges
	done.Add(len(bt.BufLabels))
	for k, v := range bt.BufLabels {
		go func(k string, v bLabel) {
			defer done.Done()
			v.nesting = append(v.nesting, "labels", k)
			v.mapConnections(c)
		}(k, v)
	}
	done.Wait()
	// parses and transforms nodes with their properties; then "maps" already
	// transformed nodes to labels
	done.Add(len(bt.BufNodes))
	for k, v := range bt.BufNodes {
		go func(k string, v bNode) {
			defer done.Done()
			v.nesting = append(v.nesting, "nodes", k)
			v.toActual(c)
		}(k, v)
	}
	done.Wait()
	// adds labels properties to according nodes using labels; then maps already
	// parsed nodes with each other using already parsed edges and LABELS; then
	// maps nodes with each other using ONLY already parsed edges
	done.Add(len(bt.BufNodes))
	for k, v := range bt.BufNodes {
		go func(k string, v bNode) {
			defer done.Done()
			v.nesting = append(v.nesting, "nodes", k)
			v.enrichByLabels(c)
			v.mapConnections(c)
		}(k, v)
	}
	done.Wait()

	if e := c.buildErr(); e != nil {
		return nil, e
	}
	return c.res, nil
}

// Concurrently mutates bEdge to actual template Edge-struct and inserts
// it to the c; don't interrupts on error occurences and writes them into
// the error-list within context
func (be bEdge) toActual(c *context) {
	name := be.nesting[len(be.nesting)-1]
	actual := &template.TEdge{
		Typ:   name,
		Props: make(map[string]*template.TProperty),
	}
	c.setEdge(name, actual)

	done := new(sync.WaitGroup)
	done.Add(len(be.BufProps))
	for k, v := range be.BufProps {
		go func(k string, v bProperty) {
			defer done.Done()
			v.nesting = append(be.nesting, "properties", k)
			v.toActual(c)
		}(k, v)
	}
	done.Wait()
}

// Concurrently mutates bLabel to buffer Label-struct (which then used to
// enrich Node-structs) and inserts it to the context; don't interrupts
// on error occurences and writes them into the error-list within context
func (bl bLabel) toActual(c *context) {
	name := bl.nesting[len(bl.nesting)-1]
	actual := &cLabel{
		typ:   name,
		props: make(map[string]*template.TProperty),
		nodes: make(map[string]*template.TNode),
	}
	c.setLabel(name, actual)

	done := new(sync.WaitGroup)
	done.Add(len(bl.BufProps))
	for k, v := range bl.BufProps {
		go func(k string, v bProperty) {
			defer done.Done()
			v.nesting = append(bl.nesting, "properties", k)
			v.toActual(c)
		}(k, v)
	}
	done.Wait()
}

// Concurrently maps connections between labels which stored within bLabel
// and inserts them as buffer Connection-structs (which then used to create
// additional actual Connection-structs) to the context; don't interrupts
// on error occurences and writes them into the error-list within context
//
// WARNING: dont call this func until according buffer Label-struct were
// insertedinside context (e.g. according toActual()-method were executed),
// otherwise it will panic - "invalid memory address or nil pointer
// dereference"
func (bl bLabel) mapConnections(c *context) {
	done := new(sync.WaitGroup)
	for k, vs := range bl.BufConns {
		done.Add(len(vs))
		for i, v := range vs {
			go func(k string, i int, v bConnection) {
				defer done.Done()
				v.nesting = append(bl.nesting, "connections", k, strconv.Itoa(i+1))
				v.toActual(c)
			}(k, i, v)
		}
	}
	done.Wait()
}

// Concurrently mutates bNode to actual Node-struct and inserts it to the
// context; don't interrupts on error occurences and writes them into the
// error-list within context
//
// WARNING: dont call this func until according Label-struct were inserted
// inside context (e.g. according toActual()-method were executed), otherwise
// it will panic - "invalid memory address or nil pointer dereference"
func (bn bNode) toActual(c *context) {
	name := bn.nesting[len(bn.nesting)-1]
	actual := &template.TNode{
		Typ:   name,
		Props: make(map[string]*template.TProperty),
	}
	c.setNode(name, actual)

	done := new(sync.WaitGroup)
	done.Add(len(bn.BufProps))
	for k, v := range bn.BufProps {
		go func(k string, v bProperty) {
			defer done.Done()
			v.nesting = append(bn.nesting, "properties", k)
			v.toActual(c)
		}(k, v)
	}
	done.Wait()

	done.Add(len(bn.BufLabels))
	for i, k := range bn.BufLabels {
		go func(i int, k string) {
			defer done.Done()
			if c.label(k) == nil {
				e := parseError{
					append(bn.nesting, "labels", strconv.Itoa(i+1)).String(),
					fmt.Sprintf("%q node has undefined label %q to attach it to node type", name, k),
				}
				c.appendErr(e)
			} else {
				c.setNodeLabel(name, k)
				c.setLabelNode(k, name, actual)
			}
		}(i, k)
	}
	done.Wait()
}

// Concurrently enriches according template Node-struct with data from
// according tamplate Label-structs (both of them detected by bNode) -
// all within context; don't interrupts on error occurences and writes
// them into the error-list within context
//
// WARNING: dont call this func until according Label- and Node-struct
// were inserted inside context (e.g. according Node's toActual()-method
// and Label's toActual()- and mapConnections()-methods were executed),
// otherwise it will panic - "invalid memory address or nil pointer
// dereference"
func (bn bNode) enrichByLabels(c *context) {
	name := bn.nesting[len(bn.nesting)-1]
	// nil checks after the most further calls of getter-funcs will be ignored - its coorect
	// in situations when returning value should be iterated or cause "nil pointer" panic
	actual := c.node(name)
	labels := c.nodeLabels(name)

	done := new(sync.WaitGroup)
	done.Add(len(labels))
	for _, mk := range labels {
		go func(actual *template.TNode, name, mk string) {
			defer done.Done()
			// merges labels properties with main node properties
			for k, v := range c.label(mk).props {
				if c.nodeProp(name, k) == nil {
					// protects nodes properties from overwriting by labels properties
					c.setNodeProp(name, k, v)
				}
			}

			// merges labels subject nodes and edges from ALL label connections with main node
			for sk := range c.labelConnsByMain(mk) {
				for ek := range c.labelConnsByMainSubj(mk, sk) {
					conn := c.labelConn(mk, sk, ek)
					m, e := actual.Typ, conn.edge.Typ
					for _, sNode := range conn.subj.nodes {
						s := sNode.Typ
						if c.nodeConn(m, s, e) == nil {
							// protects nodes connections from overwriting by labels connections
							nConn := &template.TConnection{
								Main: actual,
								Edge: conn.edge,
								Subj: sNode,
								Min:  conn.min,
								Max:  conn.max,
							}
							c.setNodeConn(m, s, e, nConn)
						}
					}
				}
			}
		}(actual, name, mk)
	}
	done.Wait()
}

// Concurrently maps connections between template Node-structs which
// stored within bNode and inserts them as tamplate Connection-structs
// to the context; don't interrupts on error occurences and writes them
// into the error-list within context
//
// WARNING: dont call this func until according Node-struct were inserted
// inside context (e.g. according toActual()-method were executed), otherwise
// it will panic - "invalid memory address or nil pointer dereference"
func (bn bNode) mapConnections(c *context) {
	done := new(sync.WaitGroup)
	for k, vs := range bn.BufConns {
		done.Add(len(vs))
		for i, v := range vs {
			go func(k string, i int, v bConnection) {
				defer done.Done()
				v.nesting = append(bn.nesting, "connections", k, strconv.Itoa(i+1))
				v.toActual(c)
			}(k, i, v)
		}
	}
	done.Wait()
}

// Mutates bProperty to actual template Property-struct and inserts it
// to the context; don't interrupts on error occurences and writes them
// into the error-list within context
//
// WARNING: dont call this func until according Label-, Node- or Edge-struct
// were inserted inside context (e.g. according toActual()-method were
// executed), otherwise it will panic - "invalid memory address or nil
// pointer dereference"
func (bp bProperty) toActual(c *context) {
	entityType := bp.nesting[len(bp.nesting)-4]
	entity := bp.nesting[len(bp.nesting)-3]
	name := bp.nesting[len(bp.nesting)-1]
	actual := &template.TProperty{
		Key: name,
		ValRestrs: make([]*template.TRestriction, 0,
			len(bp.BufRestrs.BufValueRestr)+len(bp.BufRestrs.BufRegexpRestr),
		),
		KeyRestrs: make([]*template.TRestriction, 0,
			len(bp.BufRestrs.BufKeyValueRestr)+len(bp.BufRestrs.BufKeyRegexpRestr),
		),
	}
	switch entityType {
	case "nodes":
		c.setNodeProp(entity, name, actual)
	case "edges":
		c.setEdgeProp(entity, name, actual)
	case "labels":
		c.setLabelProp(entity, name, actual)
	}

	typs, err := toDataType(bp.BufType)
	if err != nil {
		e := parseError{
			append(bp.nesting, "type").String(),
			err.Error(),
		}
		c.appendErr(e)
	}
	actual.Typ = typs.T
	if v := typs.Vt; v != template.TNull {
		actual.ValTyp = v
	}
	if v := typs.Kt; v != template.TNull {
		actual.KeyTyp = v
	}

	bp.BufRestrs.nesting = append(bp.nesting, "restrictions")
	bp.BufRestrs.toActual(c)
}

// Mutates bRestrictions to actual template Restriction-struct and
// inserts it to the context; don't interrupts on error occurences
// and writes them into the error-list within context
//
// WARNING: dont call this func until according Property-struct were
// inserted inside context (e.g. according toActual()-method were
// executed), otherwise it will panic - "invalid memory address or nil
// pointer dereference"
func (br bRestrictions) toActual(c *context) {
	entityType := br.nesting[len(br.nesting)-5]
	entity := br.nesting[len(br.nesting)-4]
	prop := br.nesting[len(br.nesting)-2]

	var p *template.TProperty
	var typs typeBuffer
	switch entityType {
	case "nodes":
		entityType = "node"
		p = c.nodeProp(entity, prop)
	case "edges":
		entityType = "edge"
		p = c.edgeProp(entity, prop)
	case "labels":
		entityType = "label"
		p = c.labelProp(entity, prop)
	}
	if p == nil {
		e := parseError{
			br.nesting.String(),
			fmt.Sprintf("%s %q has undefined property %q to write restriction to it", entity, entityType, prop),
		}
		c.appendErr(e)
	} else {
		typs = typeBuffer{
			T:  p.Typ,
			Vt: p.ValTyp,
			Kt: p.KeyTyp,
		}
	}

	for i, v := range br.BufValueRestr {
		r, err := mutateRestr(typs, template.TValue, v)
		if err != nil {
			e := parseError{
				append(br.nesting, "values", strconv.Itoa(i+1)).String(),
				err.Error(),
			}
			c.appendErr(e)
		}
		p.ValRestrs = append(p.ValRestrs, r)
	}
	for i, v := range br.BufRegexpRestr {
		r, err := mutateRestr(typs, template.TRegExp, v)
		if err != nil {
			e := parseError{
				append(br.nesting, "regexps", strconv.Itoa(i+1)).String(),
				err.Error(),
			}
			c.appendErr(e)
		}
		p.ValRestrs = append(p.ValRestrs, r)
	}

	for i, v := range br.BufKeyValueRestr {
		r, err := mutateRestr(typs, template.TKeyValue, v)
		if err != nil {
			e := parseError{
				append(br.nesting, "key_values", strconv.Itoa(i+1)).String(),
				err.Error(),
			}
			c.appendErr(e)
		}
		p.KeyRestrs = append(p.KeyRestrs, r)
	}
	for i, v := range br.BufKeyRegexpRestr {
		r, err := mutateRestr(typs, template.TKeyRegExp, v)
		if err != nil {
			e := parseError{
				append(br.nesting, "key_regexps", strconv.Itoa(i+1)).String(),
				err.Error(),
			}
			c.appendErr(e)
		}
		p.KeyRestrs = append(p.KeyRestrs, r)
	}
}

// Mutates bConnection to actual template Connection-struct and inserts
// it to the context; don't interrupts on error occurences and writes them
// into the error-list within context
//
// WARNING: dont call this func until according actual Edge- and Node-structs
// or Edge- and Labels-structs were inserted inside context (e.g. according
// toActual()-methods were executed), depending in which nesting context
// this func were called
func (bc bConnection) toActual(c *context) {
	mainType := bc.nesting[len(bc.nesting)-5]
	main := bc.nesting[len(bc.nesting)-4]
	subj := bc.nesting[len(bc.nesting)-2]
	edge := bc.BufEdge

	var m, s interface{}
	var e *template.TEdge
	var match bool
	switch mainType {
	case "nodes":
		mainType = "node"
		m = c.node(main)
		s = c.node(subj)
		match = c.nodeConn(main, subj, edge) != nil
	case "labels":
		mainType = "label"
		m = c.label(main)
		s = c.label(subj)
		match = c.labelConn(main, subj, edge) != nil
	}
	e = c.edge(edge)

	var err bool
	if s == nil {
		err = true
		e := parseError{
			bc.nesting.String(),
			fmt.Sprintf("%[1]s %[2]q has undefined subject %[1]q %[3]q to create connection", mainType, main, subj),
		}
		c.appendErr(e)
	}
	if e == nil {
		err = true
		e := parseError{
			bc.nesting.String(),
			fmt.Sprintf("%s %q has undefined edge %q to create connection", mainType, main, bc.BufEdge),
		}
		c.appendErr(e)
	}
	if match {
		err = true
		e := parseError{
			bc.nesting.String(),
			fmt.Sprintf("there is already exists such %s-connection with %q main, edge - %q and subject - %q", mainType, main, edge, subj),
		}
		c.appendErr(e)
	}
	if bc.BufRatio.Min < 0 {
		err = true
		e := parseError{
			append(bc.nesting, "ratio", "min").String(),
			"\"min\" can't be less than 0",
		}
		c.appendErr(e)
	}
	if bc.BufRatio.Max == 0 || bc.BufRatio.Max < -1 {
		err = true
		e := parseError{
			append(bc.nesting, "ratio", "max").String(),
			"\"max\" can't be equal to 0 or be less than -1 (-1 is considered as positive infinity)",
		}
		c.appendErr(e)
	}

	if !err {
		// its necessary to not insert any Connection-struct within c if any error
		// occurs cus further reusing of this struct may cause hard-to-search errors
		switch mainType {
		case "node":
			actual := &template.TConnection{
				Main: m.(*template.TNode),
				Edge: e,
				Subj: s.(*template.TNode),
				Min:  bc.BufRatio.Min,
				Max:  bc.BufRatio.Max,
			}
			c.setNodeConn(main, subj, edge, actual)
		case "label":
			actual := &cLConnection{
				main: m.(*cLabel),
				edge: e,
				subj: s.(*cLabel),
				min:  bc.BufRatio.Min,
				max:  bc.BufRatio.Max,
			}
			c.setLabelConn(main, subj, edge, actual)
		}
	}
}
