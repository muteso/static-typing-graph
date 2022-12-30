package parser

import (
	"fmt"
	"stg/template"
	"sync"
)

// Main parsing type, which contains all necessary information
// about parsing process and has some concurrent synchronisation
// tools
type context struct {
	res  *template.TemplateHolder
	ls   map[string]*cLabel                             // [label]
	nls  map[string][]string                            // [node] -> labels names
	lcn  map[string]map[string]map[string]*cLConnection // [main label] [subj label] [edge]
	errs []parseError
	*sync.Mutex
}

// Context label type - contains type name, nodes (in which label
// properties and connections are included) and properties
type cLabel struct {
	typ   string
	props map[string]*template.TProperty
	nodes map[string]*template.TNode
}

// Context connection type - represents bound between main node with
// specified label and subject node with specified label, connected by
// Edge, which ALWAYS directed from main node to subject; contains main
// node, edge, subject node and minimum and maximum possible amount of
// connections between ONE unique main node and ANY amount of unique
// subject nodes using edge
type cLConnection struct {
	main *cLabel
	edge *template.TEdge
	subj *cLabel

	min int
	max int
}

// Creates and returns new context-struct
func newContext() *context {
	c := &context{
		res: &template.TemplateHolder{
			Nodes: make(map[string]*template.TNode),
			Edges: make(map[string]*template.TEdge),
			Conns: make(map[string]map[string]map[string]*template.TConnection),
		},
		ls:    make(map[string]*cLabel),
		nls:   make(map[string][]string),
		lcn:   make(map[string]map[string]map[string]*cLConnection),
		errs:  make([]parseError, 0),
		Mutex: new(sync.Mutex),
	}
	return c
}

// ---------------------- ERRORS ---------------------- //

// Appends e to list of occured errors
func (c *context) appendErr(e parseError) {
	c.errs = append(c.errs, e)
}

// Builds and returns error which consists of all occured
// errors during parsing; returns nil if no errors occured
func (c context) buildErr() error {
	if len(c.errs) == 0 {
		return nil
	}
	res := ""

	for _, v := range c.errs {
		res += v.Error() + "\n"
	}

	return fmt.Errorf("%s", res)
}

// ---------------------- GETTERS ---------------------- //

// Returns Node-struct (with n type name) at success; returns nil
// otherwise
//
// Concurrent safe
func (c context) node(n string) *template.TNode {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		return nil
	}
	v := c.res.Nodes[n]
	return v
}

// Returns Property-struct with p property key name of node with n
// type name at success; returns nil otherwise
//
// Concurrent safe
func (c context) nodeProp(n, p string) *template.TProperty {
	c.Lock()
	defer c.Unlock()
	if c.res == nil || c.res.Nodes[n] == nil {
		return nil
	}
	v := c.res.Nodes[n].Props[p]
	return v
}

// Returns slice of label names of node with n type name at success;
// returns nil otherwise
func (c context) nodeLabels(n string) []string {
	c.Lock()
	defer c.Unlock()
	v := c.nls[n]
	return v
}

// Returns Label-struct (with l type name) at success; returns nil
// otherwise
//
// Concurrent safe
func (c context) label(l string) *cLabel {
	c.Lock()
	defer c.Unlock()
	v := c.ls[l]
	return v
}

// Returns Property-struct with p property key name of label with l
// type name at success; returns nil otherwise
//
// Concurrent safe
func (c context) labelProp(l, p string) *template.TProperty {
	c.Lock()
	defer c.Unlock()
	if c.ls[l] == nil {
		return nil
	}
	v := c.ls[l].props[p]
	return v
}

// Returns map of nodes (node names mapped to Node-structs) of label
// with l type name at success; returns nil otherwise
//
// Concurrent safe
func (c context) labelNodes(l string) map[string]*template.TNode {
	c.Lock()
	defer c.Unlock()
	if c.ls[l] == nil {
		return nil
	}
	v := c.ls[l].nodes
	return v
}

// Returns Edge-struct (with e type name) at success; returns nil
// otherwise
//
// Concurrent safe
func (c context) edge(e string) *template.TEdge {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		return nil
	}
	v := c.res.Edges[e]
	return v
}

// Returns Property-struct with p property key name of edge with e
// type name at success; returns nil otherwise
//
// Concurrent safe
func (c context) edgeProp(e, p string) *template.TProperty {
	c.Lock()
	defer c.Unlock()
	if c.res == nil || c.res.Edges[e] == nil {
		return nil
	}
	v := c.res.Edges[e].Props[p]
	return v
}

// Returns map of Connection-structs (which "origins" from label with
// m type name) at success; returns nil otherwise
//
// Concurrent safe
func (c context) labelConnsByMain(m string) map[string]map[string]*cLConnection {
	c.Lock()
	defer c.Unlock()
	v := c.lcn[m]
	return v
}

// Returns map of Connection-structs (which "origins" from label with
// m type name and "ends" at label with s type name) at success; returns
// nil otherwise
//
// Concurrent safe
func (c context) labelConnsByMainSubj(m, s string) map[string]*cLConnection {
	c.Lock()
	defer c.Unlock()
	v := c.lcn[m][s]
	return v
}

// Returns Connection-struct (which "origins" from label with m type name,
// "ends" at label with s type name and has edge with e type name) at
// success; returns nil otherwise
//
// Concurrent safe
func (c context) labelConn(m, s, e string) *cLConnection {
	c.Lock()
	defer c.Unlock()
	v := c.lcn[m][s][e]
	return v
}

// Returns map of Connection-structs (which "origins" from node with m
// type name) at success; returns nil otherwise
//
// Concurrent safe
func (c context) nodeConnsByMain(m string) map[string]map[string]*template.TConnection {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		return nil
	}
	v := c.res.Conns[m]
	return v
}

// Returns map of Connection-structs (which "origins" from node with m
// type name and "ends" at node with s type name) at success; returns
// nil otherwise
//
// Concurrent safe
func (c context) nodeConnsByMainSubj(m, s string) map[string]*template.TConnection {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		return nil
	}
	v := c.res.Conns[m][s]
	return v
}

// Returns Connection-struct (which "origins" from node with m type
// name, "ends" at node with s type name and has edge with e type name)
// at success; returns nil otherwise
//
// Concurrent safe
func (c context) nodeConn(m, s, e string) *template.TConnection {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		return nil
	}
	v := c.res.Conns[m][s][e]
	return v
}

// ---------------------- SETTERS ---------------------- //

// Sets node Node-struct with n type name within context
//
// Concurrent safe
func (c *context) setNode(n string, node *template.TNode) {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		c.res = &template.TemplateHolder{}
	}
	if c.res.Nodes == nil {
		c.res.Nodes = make(map[string]*template.TNode)
	}
	c.res.Nodes[n] = node
}

// Sets prop Property-struct with p property key name within node (with
// n type name) within context
//
// # Concurrent safe
//
// If relevant underlying data is absent - allocates it
func (c *context) setNodeProp(n, p string, prop *template.TProperty) {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		c.res = &template.TemplateHolder{}
	}
	if c.res.Nodes == nil {
		c.res.Nodes = make(map[string]*template.TNode)
	}
	if c.res.Nodes[n] == nil {
		c.res.Nodes[n] = &template.TNode{
			Typ: n,
		}
	}
	if c.res.Nodes[n].Props == nil {
		c.res.Nodes[n].Props = make(map[string]*template.TProperty)
	}
	c.res.Nodes[n].Props[p] = prop
}

// Appends l label type name within node (with n type name) within context
//
// # Concurrent safe
//
// If relevant underlying data is absent - allocates it
func (c *context) setNodeLabel(n, l string) {
	c.Lock()
	defer c.Unlock()
	if c.nls == nil {
		c.nls = make(map[string][]string)
	}
	if c.nls[n] == nil {
		c.nls[n] = make([]string, 0, 1)
	}
	c.nls[n] = append(c.nls[n], l)
}

// Sets label Label-struct with l label name within context
//
// # Concurrent safe
//
// If relevant underlying data is absent - allocates it
func (c *context) setLabel(l string, label *cLabel) {
	c.Lock()
	defer c.Unlock()
	if c.ls == nil {
		c.ls = make(map[string]*cLabel)
	}
	c.ls[l] = label
}

// Sets prop Property-struct with p property key name within label (with l
// type name) within context
//
// # Concurrent safe
//
// If relevant underlying data is absent - allocates it
func (c *context) setLabelProp(l, p string, prop *template.TProperty) {
	c.Lock()
	defer c.Unlock()
	if c.ls == nil {
		c.ls = make(map[string]*cLabel)
	}
	if c.ls[l] == nil {
		c.ls[l] = &cLabel{
			typ: l,
		}
	}
	if c.ls[l].props == nil {
		c.ls[l].props = make(map[string]*template.TProperty)
	}
	c.ls[l].props[p] = prop
}

// Sets node Node-struct with n type name within label (with l type name)
// within context
//
// # Concurrent safe
//
// If relevant underlying data is absent - allocates it
func (c *context) setLabelNode(l, n string, node *template.TNode) {
	c.Lock()
	defer c.Unlock()
	if c.ls == nil {
		c.ls = make(map[string]*cLabel)
	}
	if c.ls[l] == nil {
		c.ls[l] = &cLabel{
			typ: l,
		}
	}
	if c.ls[l].nodes == nil {
		c.ls[l].nodes = make(map[string]*template.TNode)
	}
	c.ls[l].nodes[n] = node
}

// Sets edge Edge-struct with e type name within context
//
// # Concurrent safe
//
// If relevant underlying data is absent - allocates it
func (c *context) setEdge(e string, edge *template.TEdge) {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		c.res = &template.TemplateHolder{}
	}
	if c.res.Edges == nil {
		c.res.Edges = make(map[string]*template.TEdge)
	}
	c.res.Edges[e] = edge
}

// Sets prop Property-struct with p property key name within edge (with e
// type name) witin context
//
// # Concurrent safe
//
// If relevant underlying data is absent - allocates it
func (c *context) setEdgeProp(e, p string, prop *template.TProperty) {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		c.res = &template.TemplateHolder{}
	}
	if c.res.Edges == nil {
		c.res.Edges = make(map[string]*template.TEdge)
	}
	if c.res.Edges[e] == nil {
		c.res.Edges[e] = &template.TEdge{
			Typ: e,
		}
	}
	if c.res.Edges[e].Props == nil {
		c.res.Edges[e].Props = make(map[string]*template.TProperty)
	}
	c.res.Edges[e].Props[p] = prop
}

// // Sets cs map of Connection-structs (which "origins" from node with m
// // type name) within context
// //
// // # Concurrent safe
// //
// // If relevant underlying map is absent - allocates it
// func (c *context) setNodeConnsByMain(m string, cs map[string]map[string]*template.TConnection) {
// 	c.Lock()
// 	defer c.Unlock()
// 	c.res.Conns[m] = cs
// }

// // Sets cs map of Connection-structs (which "origins" from node with m
// // type name and "ends" at node with s type name) within context
// //
// // # Concurrent safe
// //
// // If relevant underlying map is absent - allocates it
// func (c *context) setNodeConnsByMainSubj(m, s string, cs map[string]*template.TConnection) {
// 	c.Lock()
// 	defer c.Unlock()
// 	if _, ok := c.res.Conns[m]; !ok {
// 		c.res.Conns[m] = make(map[string]map[string]*template.TConnection)
// 	}
// 	c.res.Conns[m][s] = cs
// }

// Sets conn Connection-struct (which "origins" from node with m type
// name, "ends" at node with s type name and has edge with e type name)
// within context
//
// # Concurrent safe
//
// If relevant underlying map is absent - allocates it
func (c *context) setNodeConn(m, s, e string, conn *template.TConnection) {
	c.Lock()
	defer c.Unlock()
	if c.res == nil {
		c.res = &template.TemplateHolder{}
	}
	if c.res.Conns == nil {
		c.res.Conns = make(map[string]map[string]map[string]*template.TConnection)
	}
	if c.res.Conns[m] == nil {
		c.res.Conns[m] = make(map[string]map[string]*template.TConnection)
	}
	if c.res.Conns[m][s] == nil {
		c.res.Conns[m][s] = make(map[string]*template.TConnection)
	}
	c.res.Conns[m][s][e] = conn
}

// // Sets cs map of Connection-structs (which "origins" from label with m
// // type name) within context
// //
// // # Concurrent safe
// //
// // If relevant underlying map is absent - allocates it
// func (c *context) setLabelConnsByMain(m string, cs map[string]map[string]*cLConnection) {
// 	c.Lock()
// 	defer c.Unlock()
// 	c.lcn[m] = cs
// }

// // Sets cs map of Connection-structs (which "origins" from label with m
// // type name and "ends" at label with s type name) within context
// //
// // # Concurrent safe
// //
// // If relevant underlying map is absent - allocates it
// func (c *context) setLabelConnsByMainSubj(m, s string, cs map[string]*cLConnection) {
// 	c.Lock()
// 	defer c.Unlock()
// 	if _, ok := c.lcn[m]; !ok {
// 		c.lcn[m] = make(map[string]map[string]*cLConnection)
// 	}
// 	c.lcn[m][s] = cs
// }

// Sets conn Connection-struct (which "origins" from label with m type
// name, "ends" at label with s type name and has edge with e type name)
// within context
//
// # Concurrent safe
//
// If relevant underlying map is absent - allocates it
func (c *context) setLabelConn(m, s, e string, conn *cLConnection) {
	c.Lock()
	defer c.Unlock()
	if c.lcn == nil {
		c.lcn = make(map[string]map[string]map[string]*cLConnection)
	}
	if c.lcn[m] == nil {
		c.lcn[m] = make(map[string]map[string]*cLConnection)
	}
	if c.lcn[m][s] == nil {
		c.lcn[m][s] = make(map[string]*cLConnection)
	}
	c.lcn[m][s][e] = conn
}
