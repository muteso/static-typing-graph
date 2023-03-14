package template

import (
	"reflect"
	"strings"
	"unicode"
)

// Searches Node- and Edge-structs within TemplateHolder-struct using v
// underlying data and returns them if found; otherwise returns nil
//
// WARNING: under the hood func tries to find types name of the v underlying
// data firstly within type definition and then within v properties values;
// thus this func can evaluate ONLY structs, maps and map-based custom types
func searchUnknownTypeName(t TemplateHolder, v reflect.Value) (*TNode, *TEdge) {
	typName := v.Type().Name()
	asNode := t.Nodes[typName]
	asEdge := t.Edges[typName]
	if asNode == nil && asEdge == nil {
		switch v.Kind() {
		case reflect.Map:
			for iter := v.MapRange(); iter.Next(); {
				val := iter.Value()
				if val.Kind() == reflect.Interface {
					val = val.Elem()
				}
				vTyp := val.String()

				if asNode == nil {
					asNode = t.Nodes[vTyp]
				}
				if asEdge == nil {
					asEdge = t.Edges[vTyp]
				}
			}
		case reflect.Struct:
			for i := 0; i < v.Type().NumField(); i++ {
				vTyp := v.Field(i).String()
				if asNode == nil {
					asNode = t.Nodes[vTyp]
				}
				if asEdge == nil {
					asEdge = t.Edges[vTyp]
				}
			}
		}
	}
	return asNode, asEdge
}

// Tries to validate underlying data of v as map data type using t Node-struct
// as validator and returns true and nil on success
func validateUnknownMapAsNode(t *TNode, v reflect.Value) (bool, error) {
	keys, err := evaluateUnknownMapPropertyKeys(t.Typ, v, t.Props)
	if err != nil {
		return false, err
	}

	if ok, err := validateUnknownMapProperties(t.Typ, keys, v, t.Props); !ok {
		return false, err
	}

	return true, nil
}

// Tries to validate underlying data of v as map data type using t Edge-struct
// as validator and returns true and nil on success
func validateUnknownMapAsEdge(t *TEdge, v reflect.Value) (bool, error) {
	keys, err := evaluateUnknownMapPropertyKeys(t.Typ, v, t.Props)
	if err != nil {
		return false, err
	}

	if ok, err := validateUnknownMapProperties(t.Typ, keys, v, t.Props); !ok {
		return false, err
	}

	return true, nil
}

// Tries to validate underlying data of v as struct data type using t Node-struct
// as validator and returns true and nil on success
func validateUnknownStructAsNode(t *TNode, v reflect.Value) (bool, error) {
	keys, err := evaluateUnknownStructPropertyKeys(t.Typ, v, t.Props)
	if err != nil {
		return false, err
	}

	if ok, err := validateUnknownStructProperties(t.Typ, keys, v, t.Props); !ok {
		return false, err
	}

	return true, nil
}

// Tries to validate underlying data of v as struct data type using t Edge-struct
// as validator and returns true and nil on success
func validateUnknownStructAsEdge(t *TEdge, v reflect.Value) (bool, error) {
	keys, err := evaluateUnknownStructPropertyKeys(t.Typ, v, t.Props)
	if err != nil {
		return false, err
	}

	if ok, err := validateUnknownStructProperties(t.Typ, keys, v, t.Props); !ok {
		return false, err
	}

	return true, nil
}

// Maps vKeys values properties keys with tKeys template properties keys and
// returns "valid" map of keys (where vKeys keys is mapped to tKeys values) and
// nil if match any contradictions between them
func mapKeys(tKeys, vKeys []string) map[string]string {
	validKeys := make(map[string]string)

	tempKeys := make(map[string]string)
	for _, k := range tKeys {
		tempKeys[strings.ToLower(k)] = k
	}
	for _, k := range vKeys {
		if templateKey, ok := tempKeys[strings.ToLower(k)]; ok {
			validKeys[k] = templateKey
		} else {
			return nil
		}
	}
	return validKeys
}

// Parses and compares ps properties (as origin) and v underlying data (considered
// as map) properties (as implementation of origin) and merges them into a single
// "valid" map (where v keys is mapped to ps keys) of keys which is returned on
// success; omit argument is used when v underlying data contains its own type name
// as one the properties value and thus should be omitted - otherwise it can be ""
func evaluateUnknownMapPropertyKeys(omit string, v reflect.Value, ps map[string]*TProperty) (map[string]string, error) {
	vKeys := make([]string, 0)
	for iter := v.MapRange(); iter.Next(); {
		val := iter.Value()
		if val.Kind() == reflect.Interface {
			val = val.Elem()
		}
		if omit == val.String() {
			continue
		}
		vKeys = append(vKeys, iter.Key().String())
	}
	tKeys := make([]string, 0, len(ps))
	for k := range ps {
		tKeys = append(tKeys, k)
	}
	if err := comparePropertyKeys(tKeys, vKeys); err != nil {
		return nil, err
	}
	return mapKeys(tKeys, vKeys), nil
}

// Parses and compares ps properties (as origin) and v underlying data (considered
// as struct) properties (as implementation of origin) and merges them into a single
// "valid" map (where v keys is mapped to ps keys) of keys which is returned on success;
// omit argument is used when v underlying data contains its own type name as one the
// properties value and thus should be omitted - otherwise it can be ""
func evaluateUnknownStructPropertyKeys(omit string, v reflect.Value, ps map[string]*TProperty) (map[string]string, error) {
	vKeys := make([]string, 0)
	vUpperKeys := make([]string, 0)
	vLowerKeys := make([]string, 0)
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		if omit == v.Field(i).String() {
			continue
		}
		var uk, lk string
		for i, dig := range vt.Field(i).Name {
			udig, ldig := dig, dig
			if i == 0 {
				udig = unicode.ToUpper(dig)
				ldig = unicode.ToLower(dig)
			}
			uk += string(udig)
			lk += string(ldig)
		}
		vKeys = append(vKeys, vt.Field(i).Name)
		vUpperKeys = append(vUpperKeys, uk)
		vLowerKeys = append(vLowerKeys, lk)
	}
	tKeys := make([]string, 0, len(ps))
	for k := range ps {
		tKeys = append(tKeys, k)
	}
	if err := comparePropertyKeys(tKeys, vUpperKeys); err != nil {
		if err := comparePropertyKeys(tKeys, vLowerKeys); err != nil {
			return nil, err
		}
	}

	return mapKeys(tKeys, vKeys), nil
}

// Tries to validate v underlying data as map data type using ks "valid" keys and ps as
// validators and returns true and nil on success; omit argument is used when v underlying
// data contains its own type name as one the properties value and thus should be omitted -
// otherwise it can be ""
func validateUnknownMapProperties(omit string, ks map[string]string, v reflect.Value, ps map[string]*TProperty) (bool, error) {
	for vk, tk := range ks {
		p := v.MapIndex(reflect.ValueOf(vk))
		if omit == p.String() {
			continue
		}
		tp := *ps[tk]

		ok, err := evaluateProperty(tp, p.Interface())
		if !ok {
			return ok, err
		}
	}
	return true, nil
}

// Tries to validate v underlying data as struct data type using ks "valid" keys and ps as
// validators and returns true and nil on success; omit argument is used when v underlying
// data contains its own type name as one the properties value and thus should be omitted -
// otherwise it can be ""
func validateUnknownStructProperties(omit string, ks map[string]string, v reflect.Value, ps map[string]*TProperty) (bool, error) {
	for vk, tk := range ks {
		p := v.FieldByName(vk)
		if omit == p.String() {
			continue
		}
		tp := *ps[tk]

		ok, err := evaluateProperty(tp, p.Interface())
		if !ok {
			return ok, err
		}
	}
	return true, nil
}
