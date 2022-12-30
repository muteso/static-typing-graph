package validation

// Validator interface that is used as validation tool - it's necessary
// to correctly implement ValidateNode-, ValidateEdge- and ValidateTriplet-funcs;
// contrary of this ValidateUnknown-func should be used as "last chance"-tool
// in situations where the given validating argument is unkhown and doesn't
// implement any of the Node-, Edge- and Triplet-interfaces; because
// of this, if you want to implement Validator-interface by your own
// and don't want to give to the user such "last chance"-tool you can
// implement ValidateUnknown-func as a stub, for example like this:
//
// 	func (vr YourOwnType) ValidateUnknown(v interface{}) (bool, error) {
//		return false, errors.New("given value doesn't implement any of the Node-, Edge- or Triplet-interfaces")
// 	}
type Validator interface {
	ValidateNode(Node) (bool, error)
	ValidateEdge(Edge) (bool, error)
	ValidateTriplet(Triplet) (bool, error)
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

// Buffer type that implements Node- and Edge-interface and thus can be
// used to mutate any given data to those interfaces values
type entityHolder struct {
	typ   string
	props map[string]interface{}
}

// Returns type name of entity like it is considered as node
func (h entityHolder) GetNodeType() string {
	return h.typ
}

// Returns type name of edge like it is considered as edge
func (h entityHolder) GetEdgeType() string {
	return h.typ
}

// Returns slice of property keys of entity; returns empty slice if there
// is no properties
func (h entityHolder) GetKeys() []string {
	res := make([]string, 0, len(h.props))
	for k := range h.props {
		res = append(res, k)
	}
	return res
}

// Returns value of property by key from entity; returns nil if such key
// is absent
func (h entityHolder) GetProp(key string) (interface{}, bool) {
	v, ok := h.props[key]
	return v, ok
}

// Triplet interface that uses RDF-representation as sematic reference;
// any type that can return main entity, subject entity and entity, that
// semanticly connects them, implements that interface
type Triplet interface {
	Main() Node
	Edge() Edge
	Subj() Node
}

// Buffer type that implements Triplet interface and thus can be used to
// mutate 3 given Entity-interface values to Triplet-interface value
type tripletHolder struct {
	m Node
	e Edge
	s Node
}

// Returns main Entity-interface value
func (t tripletHolder) Main() Node {
	return t.m
}

// Returns edge Entity-interface value
func (t tripletHolder) Edge() Edge {
	return t.e
}

// Returns subject Entity-interface value
func (t tripletHolder) Subj() Node {
	return t.s
}

// // refactor
// type Graph interface {
//	GetNode() error
// 	AddNode() error
// 	UpdateNode() error
// 	DeleteNode() error
// }

// type graphHolder struct {
// 	graph []*graphNode
// 	pass  []*graphNode
// }

// type graphNode struct {
// 	parentEdge *Edge
// 	node       *Node
// 	childs     []*Node
// }

// func (g graphHolder) GetNode() error {
// 	return nil
// }

// func (g graphHolder) AddNode() error {
// 	return nil
// }

// func (g graphHolder) ConnectNode() error {
// 	return nil
// }

// func (g graphHolder) UpdateNode() error {
// 	return nil
// }

// func (g graphHolder) DeleteNode() error {
// 	return nil
// }
