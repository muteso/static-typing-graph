package validation

// Validator interface that is used as validation tool - it's necessary
// to correctly implement ValidateNode-, ValidateEdge-, ValidateTriplet-
// and ValidateGraph-funcs; contrary of this ValidateUnknown-func should
// be used as "last chance"-tool in situations where the given validating
// argument is unkhown and doesn't implement any of the Node-, Edge-, Graph-
// or Triplet-interfaces; because of this, if you want to implement
// Validator-interface by your own and don't want to give to the user such
// "last chance"-tool you can implement ValidateUnknown-func as a stub, for
// example like this:
//
//	func (vr YourOwnType) ValidateUnknown(v interface{}) (bool, error) {
//		return false, errors.New("given value doesn't implement any of the Node-, Edge-, Graph- or Triplet-interfaces")
//	}
//
// Recommendation: t's better to implement ValidateUnknown-func ONLY for
// recognizing SINGLE graph units, i.e. nodes or edges, but not for triplets
// or graphs
type Validator interface {
	ValidateNode(Node) (bool, error)
	ValidateEdge(Edge) (bool, error)
	ValidateTriplet(Triplet) (bool, error)
	ValidateDuplet(Duplet) (bool, error)
	ValidateGraph(Graph) (bool, error)
	ValidateUnknown(interface{}) (bool, error)
}

// Node interface - represents graph node which can return his type's
// name and properties
type Node interface {
	GetNodeType() string
	GetKeys() []string
	GetProp(string) (interface{}, bool)
}

// Edge interface - represents graph edge which can return his type's
// name and properties
type Edge interface {
	GetEdgeType() string
	GetKeys() []string
	GetProp(string) (interface{}, bool)
}

// Triplet interface that uses RDF-representation as semantic reference;
// any type that can return main entity, subject entity and entity, that
// semanticly connects them, implements that interface
type Triplet interface {
	Main() Node
	Edge() Edge
	Subj() Node
}

// Duplet interface that uses "halfed" RDF-representation as semantic
// reference - Duplet-interface express only "tail" or "head" of Triplet-interface
// (without main entity or without subject entity accordingly); any type
// that can return "head"/"tail" entity and entity, that semanticly connects
// "head"/"tail" entity with "invisible" "tail"/"head" entity, implements
// that interface
type Duplet interface {
	Edge() Edge
	Node() Node
}

// Graph interface - represents some amount of nodes interconnected (or
// not) with ecah other by edges; have methods for both introspection and
// mutation
//
// WARNING: Graph-interface should NOT assign some sort of ID's to the nodes
// and edges and should NOT do some sort of indexing, so it can't be used as
// fully functional inmemory DB; also, because of this Graph-interface should
// NOT allow to contain identical nodes and edges within implementing data type
type Graph interface {
	// "Deconstructs" whole graph to triplets and returns them; does NOT actually
	// mutate any data within graph - it is just introspection operation
	GetTriplets() []Triplet

	// Returns ALL unique nodes from graph
	GetNodes() []Node
	// Returns ALL childs (subject nodes) of given node in form of duplets (edge
	// and subject node)
	GetNodeChilds(Node) []Duplet

	// Returns ALL unique nodes from graph filtered by given main node type-name
	GetNodesByType(string) []Node
	// Returns ALL possible triplets where the given node is main node
	GetTripletsByNode(Node) []Triplet
	// Returns ALL possible triplets filtered by given main node type-name (required),
	// subject node type-name (unrequired) and edge type-name (unrequired); the order
	// of arguments should be the exactly as it enumerated previously
	GetTripletsByType(string, ...string) []Triplet

	// Inserts a single node to the graph
	AddNode(Node) bool
	// Removes a single node from the graph (with all its interconnections)
	RemoveNode(Node) bool
	// Inserts all triplet entities to the graph, but if some of them are already
	// presented in graph - inserts ONLY missed entities (or omits insertation in case
	// of already presented edge)
	AddTriplet(Triplet) bool
	// Removes a single edge from the graph (which lies between exact main node and
	// exact subject node)
	RemoveEdge(Triplet) bool
}
