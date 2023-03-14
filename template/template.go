package template

import (
	"fmt"
	"reflect"
	"stg/validation"
)

// Main type that contains all other template-types and their
// interconnections: Nodes and Conns uses Node-type's names as
// keys, Edges uses Edge-type's names as keys and Labels uses
// Labels-type's names as keys accordingly
type TemplateHolder struct {
	Nodes map[string]*TNode                             // [node]
	Edges map[string]*TEdge                             // [edge]
	Conns map[string]map[string]map[string]*TConnection // [main node][subj node][edge]
}

// for auto check of interface implementation
var _ validation.Validator = TemplateHolder{}

// Tries to validate underlying data of n as node and returns true and
// nil on success
func (t TemplateHolder) ValidateNode(n validation.Node) (bool, error) {
	typ := n.GetNodeType()
	node, ok := t.Nodes[typ]
	if !ok {
		return false, fmt.Errorf("%q-node: there is no such node type in template", typ)
	}

	vKeys := n.GetKeys()
	pKeys := make([]string, 0, len(node.Props))
	for k := range node.Props {
		pKeys = append(pKeys, k)
	}
	if err := comparePropertyKeys(pKeys, vKeys); err != nil {
		return false, fmt.Errorf("%q-node: %s", typ, err.Error())
	}

	for _, k := range vKeys {
		tp := node.Props[k]
		p, _ := n.GetProp(k)

		ok, err := evaluateProperty(*tp, p)
		if !ok {
			strErr := fmt.Sprintf("%q-node: %s", typ, err.Error())
			return ok, fmt.Errorf(strErr)
		}
	}
	return true, nil
}

// Tries to validate underlying data of e as edge and returns true and
// nil on success
func (t TemplateHolder) ValidateEdge(e validation.Edge) (bool, error) {
	typ := e.GetEdgeType()
	edge, ok := t.Edges[typ]
	if !ok {
		return false, fmt.Errorf("%q-edge: there is no such edge type in template", typ)
	}

	vKeys := e.GetKeys()
	pKeys := make([]string, 0, len(edge.Props))
	for k := range edge.Props {
		pKeys = append(pKeys, k)
	}
	if err := comparePropertyKeys(pKeys, vKeys); err != nil {
		return false, fmt.Errorf("%q-edge: %s", typ, err.Error())
	}

	for _, k := range vKeys {
		tp := edge.Props[k]
		p, _ := e.GetProp(k)

		ok, err := evaluateProperty(*tp, p)
		if !ok {
			strErr := fmt.Sprintf("%q-edge: %s", typ, err.Error())
			return ok, fmt.Errorf(strErr)
		}
	}
	return true, nil
}

// Validates underlying data of tr as node (main or subject) and 1 edge
// and returns true and nil on success
func (t TemplateHolder) ValidateDuplet(du validation.Duplet) (bool, error) {
	ok := false
	for sk := range t.Conns[du.Node().GetNodeType()] {
		_, ok = t.Conns[du.Node().GetNodeType()][sk][du.Edge().GetEdgeType()]
		if ok {
			break
		}
	}
	for mk := range t.Conns {
		_, ok = t.Conns[mk][du.Node().GetNodeType()][du.Edge().GetEdgeType()]
		if ok {
			break
		}
	}
	if !ok {
		return ok, fmt.Errorf("Duplet: such connection is not exist")
	}
	if ok, err := t.ValidateNode(du.Node()); !ok {
		return ok, fmt.Errorf("Node of duplet: " + err.Error())
	}
	if ok, err := t.ValidateEdge(du.Edge()); !ok {
		return ok, fmt.Errorf("Edge of duplet: " + err.Error())
	}
	return true, nil
}

// Validates underlying data of tr as 2 nodes (main and subject) and 1
// edge and returns true and nil on success
func (t TemplateHolder) ValidateTriplet(tr validation.Triplet) (bool, error) {
	_, ok := t.Conns[tr.Main().GetNodeType()][tr.Subj().GetNodeType()][tr.Edge().GetEdgeType()]
	if !ok {
		return ok, fmt.Errorf("Triplet: such connection is not exist")
	}
	if ok, err := t.ValidateNode(tr.Main()); !ok {
		return ok, fmt.Errorf("Main node of triplet: " + err.Error())
	}
	if ok, err := t.ValidateNode(tr.Subj()); !ok {
		return ok, fmt.Errorf("Subject node of triplet: " + err.Error())
	}
	if ok, err := t.ValidateEdge(tr.Edge()); !ok {
		return ok, fmt.Errorf("Edge of triplet: " + err.Error())
	}
	return true, nil
}

// Validates underlying data of gr as fully functional graph consisted of
// nodes and edges; returns true and nil on success
func (t TemplateHolder) ValidateGraph(gr validation.Graph) (bool, error) {
	nodes := gr.GetNodes()
	for _, mNode := range nodes {
		if ok, err := t.ValidateNode(mNode); !ok {
			return ok, fmt.Errorf("Graph: " + err.Error())
		}

		// subject nodes and edges validation
		connCount := make(map[string]map[string]int)
		childs := gr.GetNodeChilds(mNode)
		for _, child := range childs {
			edgeType := child.Edge().GetEdgeType()
			sNodeType := child.Node().GetNodeType()
			if ok, err := t.ValidateEdge(child.Edge()); !ok {
				return ok, fmt.Errorf("Graph: " + err.Error())
			}
			if ok, err := t.ValidateNode(child.Node()); !ok {
				return ok, fmt.Errorf("Graph: " + err.Error())
			}
			if connCount[sNodeType] == nil {
				connCount[sNodeType] = make(map[string]int)
			}
			connCount[sNodeType][edgeType] += 1
		}

		// connections validation
		triplets := gr.GetTripletsByNode(mNode)
		for _, tr := range triplets {
			mTyp := tr.Main().GetNodeType()
			sTyp := tr.Subj().GetNodeType()
			eTyp := tr.Edge().GetEdgeType()
			conn, ok := t.Conns[mTyp][sTyp][eTyp]
			if !ok {
				return false, fmt.Errorf(
					"Graph: %q: has wrong connection through %q-edge to %q-node",
					mTyp, eTyp, sTyp)
			}
			graphConnCount := connCount[sTyp][eTyp] // may be 0 if is not presented in map
			if conn.Min > graphConnCount {
				return false, fmt.Errorf(
					"Graph: %q: has %d connection through %q-edge to %q-node, which is less then %d minimum",
					mNode.GetNodeType(),
					graphConnCount,
					eTyp, sTyp, conn.Min)
			}
			if conn.Max < graphConnCount && conn.Max != -1 {
				return false, fmt.Errorf(
					"Graph: %q: has %d connection through %q-edge to %q-node, which is larger then %d maximum",
					mNode.GetNodeType(),
					graphConnCount,
					eTyp, sTyp, conn.Max)
			}
		}
	}
	return true, nil
}

// Tries to validate underlying data of v as node or edge and returns true
// and non-nil error on success - in this case error contains description
// about what were exactly validated - node or/and edge; in case when both
// validations (as node and as edge) were unsuccessful func returns false
// and non-nil error
//
// WARNING: under the hood func tries to find types name of the v underlying
// data firstly within type definition and then within v properties values;
// thus this func can evaluate ONLY structs, maps and map-based custom types;
// in case when func tries to validate unknown struct be care of exported and
// unexported fields - it requires that either ALL property keys in according
// (to validated struct) template definition should be capitalized or none of
// them
func (t TemplateHolder) ValidateUnknown(v interface{}) (bool, error) {
	okNode, okEdge := false, false
	errNode, errEdge := fmt.Errorf("there is no such node type in template"), fmt.Errorf("there is no such edge type in template")
	val := reflect.ValueOf(v)
	if asNode, asEdge := searchUnknownTypeName(t, val); asNode != nil || asEdge != nil {
		switch val.Kind() {
		case reflect.Map:
			if asNode != nil {
				okNode, errNode = validateUnknownMapAsNode(asNode, val)
			}
			if asEdge != nil {
				okEdge, errEdge = validateUnknownMapAsEdge(asEdge, val)
			}
		case reflect.Struct:
			if asNode != nil {
				okNode, errNode = validateUnknownStructAsNode(asNode, val)
			}
			if asEdge != nil {
				okEdge, errEdge = validateUnknownStructAsEdge(asEdge, val)
			}
		default:
			return false, fmt.Errorf("unknown value: value can't be evaluated - it's not a struct, map or map-based type")
		}

		if !okNode && !okEdge {
			return false, fmt.Errorf("unknown value: value had failed both validations: as node - %s; as edge - %s", errNode.Error(), errEdge.Error())
		} else if okNode {
			return true, fmt.Errorf("unknown value: value had been saccessfully validated as node")
		} else if okEdge {
			return true, fmt.Errorf("unknown value: value had been saccessfully validated as edge")
		} else if okNode && okEdge {
			return true, fmt.Errorf("unknown value: value had been both saccessfully validated: and as node, and as edge")
		}
	}
	return false, fmt.Errorf("unknown value: values type name doesn't match any of the templates nodes or edges")
}
