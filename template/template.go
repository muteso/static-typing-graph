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
	if _, err := evaluatePropertyKeys(pKeys, vKeys); err != nil {
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
	if _, err := evaluatePropertyKeys(pKeys, vKeys); err != nil {
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

// Validates underlying structs of tr as 2 nodes (main and subject) and 1
// edge and returns true and nil on success
func (t TemplateHolder) ValidateTriplet(tr validation.Triplet) (bool, error) {
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

// Tries to validate underlying data of v as node or edge and returns true
// and non-nil error on success - in this case error contains description
// about what were exactly validated - node or/and edge; in case when both
// validations (as node and as edge) were unsuccessful func returns false
// and non-nil error
//
// WARNING: under the hood func tries to find types name of the v underlying
// data firstly within type definition and then within v properties values;
// thus this func can evaluate ONLY structs, maps and map-based custom types
func (t TemplateHolder) ValidateUnknown(v interface{}) (bool, error) {
	okNode, okEdge := false, false
	errNode, errEdge := fmt.Errorf("there is no such node type in template"), fmt.Errorf("there is no such edge type in template")
	val := reflect.ValueOf(v)
	if asNode, asEdge := t.searchUnknownTypeName(val); asNode != nil || asEdge != nil {
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
