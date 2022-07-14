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

	v_keys := n.GetKeys()
	p_keys := make([]string, 0)
	for k := range node.Props {
		p_keys = append(p_keys, k)
	}
	if _, err := evaluatePropertyKeys(p_keys, v_keys); err != nil {
		return false, fmt.Errorf("%q-node: %s", typ, err.Error())
	}

	for _, k := range v_keys {
		tp := node.Props[k]
		p, _ := n.GetProp(k)

		ok, err := evaluateProperty(*tp, p)
		if !ok {
			str_err := fmt.Sprintf("%q-node: %s", typ, err.Error())
			return ok, fmt.Errorf(str_err)
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

	v_keys := e.GetKeys()
	p_keys := make([]string, 0)
	for k := range edge.Props {
		p_keys = append(p_keys, k)
	}
	if _, err := evaluatePropertyKeys(p_keys, v_keys); err != nil {
		return false, fmt.Errorf("%q-edge: %s", typ, err.Error())
	}

	for _, k := range v_keys {
		tp := edge.Props[k]
		p, _ := e.GetProp(k)

		ok, err := evaluateProperty(*tp, p)
		if !ok {
			str_err := fmt.Sprintf("%q-edge: %s", typ, err.Error())
			return ok, fmt.Errorf(str_err)
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
	ok_node, ok_edge := false, false
	err_node, err_edge := fmt.Errorf("there is no such node type in template"), fmt.Errorf("there is no such edge type in template")
	val := reflect.ValueOf(v)
	if as_node, as_edge := t.searchUnknownTypeName(val); as_node != nil || as_edge != nil {
		switch val.Kind() {
		case reflect.Map:
			if as_node != nil {
				ok_node, err_node = validateUnknownMapAsNode(as_node, val)
			}
			if as_edge != nil {
				ok_edge, err_edge = validateUnknownMapAsEdge(as_edge, val)
			}
		case reflect.Struct:
			if as_node != nil {
				ok_node, err_node = validateUnknownStructAsNode(as_node, val)
			}
			if as_edge != nil {
				ok_edge, err_edge = validateUnknownStructAsEdge(as_edge, val)
			}
		default:
			return false, fmt.Errorf("Unknown value: value can't be evaluated - it's not a struct, map or map-based type")
		}

		if !ok_node && !ok_edge {
			return false, fmt.Errorf("Unknown value: value had failed both validations: as node - %s; as edge - %s", err_node.Error(), err_edge.Error())
		} else if ok_node {
			return true, fmt.Errorf("Unknown value: value had been saccessfully validated as node")
		} else if ok_edge {
			return true, fmt.Errorf("Unknown value: value had been saccessfully validated as edge")
		} else if ok_node && ok_edge {
			return true, fmt.Errorf("Unknown value: value had been both saccessfully validated: and as node, and as edge")
		}
	}
	return false, fmt.Errorf("Unknown value: values type name doesn't match any of the templates nodes or edges")
}
