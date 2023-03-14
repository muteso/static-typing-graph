package validation

// Creates and returns new Node-interface value with typ type name
// and props properties
func NewNode(typ string, props map[string]interface{}) Node {
	return nodeHolder{
		typ:   typ,
		props: props,
	}
}

// Creates and returns new Edge-interface value with typ type name
// and props properties
func NewEdge(typ string, props map[string]interface{}) Edge {
	return edgeHolder{
		typ:   typ,
		props: props,
	}
}

// Creates and returns new Triplet-interface with m as main node entity,
// s as subject node entity and e as edge entity
func NewTriplet(m, s Node, e Edge) Triplet {
	return tripletHolder{
		m: m,
		e: e,
		s: s,
	}
}

// Creates and returns new Duplet-interface with n as "tail"/"head" node
// entity and e as edge entity
func NewDuplet(n Node, e Edge) Duplet {
	return dupletHolder{
		e: e,
		n: n,
	}
}

// Creates and returns new Graph-interface using ns nodes and (optionally,
// but more recommended) gr triplets for building it - this func parses any
// given node and triplet and inserts it into Graph-interface, interconnecting
// "newcoming" entities with already inserted graph-entities (using main nodes
// of triplets for infering where it should happen in case of triplets)
//
// WARNING: this func dont allow some graph design practices that leads
// to excess data storaging, as:
//   - 2 or more identical nodes (this func will create only one instance of node)
//   - 2 or more identical edges with identical directions between one unique pair
//     of nodes (this func will create only one instance of edge)
//   - identity of nodes (and edges too) is found if they have equal types and
//     ALL their properties are equal
//
// Respectively, triplets with equal main nodes will be "squashed" into a
// "single" graph node and any duplicates will be just OMITTED (i.e. identical
// triplets, nodes and edges)
//
// Also, in some cases inserting of triplets has a little bit different semanthic:
//   - already presented edge - cancels insertation
//   - already presented main node and subject node - connects them with
//     edge (adding edge)
//   - nil subject node and nil edge - inserts main node as single node
//     (adding node)
//   - nil main node and nil edge - inserts subject node as single node
//     (adding node)
//   - nil edge - inserts main node and subject node as single nodes (adding
//     nodes)
//
// Because of this it's more recommendend to use ONLY triplets to build graphs, which
// may be done by using:
//
//	graph := NewGraph(nil, tripletSlice...)
func NewGraph(ns []Node, gr ...Triplet) Graph {
	res := make(graphHolder, 0)
	for _, n := range ns {
		res.insertNode(n)
	}
	for _, tr := range gr {
		res.insertTriplet(tr)
	}
	return res
}

// Validates v (which might implements Node-, Edge-, Triplet- or Graph
// interface) underlying data using vr and returns true and nil on success;
// in cases where underlying data dont implement interfaces enumerated above
// this func still tries to validate it but in some successive cases may also
// return non-nil error explaning what was exactly validated
func Validate(vr Validator, v interface{}) (bool, error) {
	switch val := v.(type) {
	case Triplet:
		return vr.ValidateTriplet(val)
	case Duplet:
		return vr.ValidateDuplet(val)
	case Node:
		return vr.ValidateNode(val)
	case Edge:
		return vr.ValidateEdge(val)
	case Graph:
		return vr.ValidateGraph(val)
	default:
		return vr.ValidateUnknown(val)
	}
}
