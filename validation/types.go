package validation

// Buffer type that implements Node-interface and thus can be used to
// mutate any given data to those interfaces values
type nodeHolder struct {
	typ   string
	props map[string]interface{}
}

// Returns type name of entity like it is considered as node
func (h nodeHolder) GetNodeType() string {
	return h.typ
}

// Returns slice of property keys of entity; returns empty slice if there
// is no properties
func (h nodeHolder) GetKeys() []string {
	res := make([]string, 0, len(h.props))
	for k := range h.props {
		res = append(res, k)
	}
	return res
}

// Returns value of property by key from entity; returns nil if such key
// is absent
func (h nodeHolder) GetProp(key string) (interface{}, bool) {
	v, ok := h.props[key]
	return v, ok
}

// Buffer type that implements Edge-interface and thus can be used to
// mutate any given data to those interfaces values
type edgeHolder struct {
	typ   string
	props map[string]interface{}
}

// Returns type name of edge like it is considered as edge
func (h edgeHolder) GetEdgeType() string {
	return h.typ
}

// Returns slice of property keys of entity; returns empty slice if there
// is no properties
func (h edgeHolder) GetKeys() []string {
	res := make([]string, 0, len(h.props))
	for k := range h.props {
		res = append(res, k)
	}
	return res
}

// Returns value of property by key from entity; returns nil if such key
// is absent
func (h edgeHolder) GetProp(key string) (interface{}, bool) {
	v, ok := h.props[key]
	return v, ok
}

// Buffer type that implements Triplet interface and thus can be used to
// mutate 3 given Entity-interface values to Triplet-interface value
type tripletHolder struct {
	m Node
	e Edge
	s Node
}

// Returns main Node-interface value
func (t tripletHolder) Main() Node {
	return t.m
}

// Returns Edge-interface value
func (t tripletHolder) Edge() Edge {
	return t.e
}

// Returns subject Node-interface value
func (t tripletHolder) Subj() Node {
	return t.s
}

// Buffer type that implements Duplet interface and thus can be used to
// mutate 2 given Entity-interface values to Duplet-interface value
type dupletHolder struct {
	e Edge
	n Node
}

// Returns Edge-interface value
func (d dupletHolder) Edge() Edge {
	return d.e
}

// Returns "tail"/"head" Node-interface value
func (d dupletHolder) Node() Node {
	return d.n
}

// Buffer type that implements Graph-interface; a little more convinient
// than slice-based type
type graphHolder map[*graphHolderNode]struct{}

// Buffer type for graphHolder-struct, which represents node within graph;
// contains node itself and all its interconnections within graph
type graphHolderNode struct {
	node   Node                        // node itself
	childs map[*graphHolderNode][]Edge // [subj] -> Edges
}

// Returns ALL unique nodes
func (g graphHolder) GetNodes() []Node {
	return g.findNodes("")
}

// Returns ALL subject nodes of given n main node in form of edge and subject
// node pairs (duplets)
func (g graphHolder) GetNodeChilds(n Node) []Duplet {
	return g.getNodeChilds(n)
}

// Returns ALL possible main node, edge and subject node threes (triplets)
func (g graphHolder) GetTriplets() []Triplet {
	blankFilters := make([]string, 0, 3)
	return g.findTriplets(blankFilters...)
}

// Returns ALL unique nodes from graph with mTyp type name
func (g graphHolder) GetNodesByType(mTyp string) []Node {
	return g.findNodes(mTyp)
}

// Returns ALL possible main node, edge and subject node threes (triplets)
// where main node is n node
func (g graphHolder) GetTripletsByNode(n Node) []Triplet {
	return g.getGroupedTriplets(n)
}

// Returns ALL possible main node, edge and subject node threes (triplets)
// where main node has mTyp type name and (optionaly!) subject node has
// (second argument) type name and (optionaly!) edge has (third argument)
// type name; any excess arguments will be omitted
func (g graphHolder) GetTripletsByType(mTyp string, typs ...string) []Triplet {
	filters := append(append(make([]string, 0, 3), mTyp), typs...)
	return g.findTriplets(filters...)
}

// Inserts a single n node to the graph
func (g graphHolder) AddNode(n Node) bool {
	return g.insertNode(n)
}

// Removes a single n node from the graph (with all its interconnections)
func (g graphHolder) RemoveNode(n Node) bool {
	return g.deleteNode(n)
}

// Inserts all entities from tr triplet to the graph, but if some of them
// are already presented in graph - inserts ONLY missed entities; in some
// cases has a little bit different semanthic:
//   - already presented edge - cancels insertation
//   - already presented main node and subject node - connects them with
//     edge (adding edge)
//   - nil subject node and nil edge - inserts main node as single node
//     (adding node)
//   - nil main node and nil edge - inserts subject node as single node
//     (adding node)
//   - nil edge - inserts main node and subject node as single nodes (adding
//     nodes)
func (g graphHolder) AddTriplet(tr Triplet) bool {
	return g.insertTriplet(tr)
}

// Removes a single edge (between exact main node and exact subject node
// enumerated in tr triplet) from the graph
func (g graphHolder) RemoveEdge(tr Triplet) bool {
	return g.deleteEdge(tr)
}
