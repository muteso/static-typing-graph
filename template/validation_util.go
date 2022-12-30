package template

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Compares vKeys values properties keys with tKeys template properties keys
// and returns "valid" map of keys (where vKeys keys is mapped to tKeys keys)
// and nil if doesn't match any contradictions between them
func evaluatePropertyKeys(tKeys []string, vKeys []string) (map[string]string, error) {
	validKeys := make(map[string]string)

	extraKeysLowcase := make(map[string]string)
	for _, k := range tKeys {
		extraKeysLowcase[strings.ToLower(k)] = k
	}
	for _, k := range vKeys {
		lowcaseTemplateKey := strings.ToLower(k)
		templateKey, ok := extraKeysLowcase[lowcaseTemplateKey]
		if !ok {
			return nil, fmt.Errorf("validated entity has extra %q property", k)
		} else {
			validKeys[k] = templateKey
			delete(extraKeysLowcase, lowcaseTemplateKey)
		}
	}
	if len(extraKeysLowcase) != 0 {
		missedProps := ""
		for _, k := range extraKeysLowcase {
			missedProps += fmt.Sprintf("%q, ", k)
		}
		missedProps = strings.TrimSuffix(missedProps, ", ")
		return nil, fmt.Errorf("validated entity doesn't has %s properties", missedProps)
	}
	return validKeys, nil
}

// Parses underlying data of p as property and validates it using tp as validator
// and returns true and nil on success
func evaluateProperty(tp TProperty, p interface{}) (bool, error) {
	switch tp.Typ {
	case TInt:
		return evaluatePropertyAsInt(tp, p)
	case TFloat:
		return evaluatePropertyAsFloat(tp, p)
	case TString:
		return evaluatePropertyAsString(tp, p)
	case TBool:
		return evaluatePropertyAsBool(tp, p)
	case TDateTime:
		return evaluatePropertyAsDateTime(tp, p)
	case TArray:
		return evaluatePropertyAsArr(tp, p)
	case TMap:
		return evaluatePropertyAsMap(tp, p)
	}
	return false, fmt.Errorf("%q-property: value doesn't match any possible data type", tp.Key)
}

// Parses underlying data of p as int data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsInt(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertInt(p)
	if !ok {
		return false, fmt.Errorf("%q-property: %q value doesn't match \"int\" data type", tp.Key, val)
	}
	for _, restr := range tp.Restrs {
		switch restr.RestrTyp {
		case TValue:
			if val != restr.Restr {
				return false, fmt.Errorf("%q-property: %q value doesn't match %q restriction", tp.Key, val, restr.Restr)
			}
		case TRegExp:
			if !restr.Restr.(*regexp.Regexp).MatchString(strconv.FormatInt(int64(val), 10)) {
				return false, fmt.Errorf("%q-property: %q value doesn't match %q regular expression", tp.Key, val, restr.Restr)
			}
		}
	}
	return true, nil
}

// Parses underlying data of p as float data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsFloat(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertFloat(p)
	if !ok {
		return false, fmt.Errorf("%q-property: value \"%f\" doesn't match \"float\" data type", tp.Key, val)
	}
	for _, restr := range tp.Restrs {
		switch restr.RestrTyp {
		case TValue:
			if val != restr.Restr {
				return false, fmt.Errorf("%q-property: \"%f\" value doesn't match %f restriction", tp.Key, val, restr.Restr)
			}
		case TRegExp:
			if !restr.Restr.(*regexp.Regexp).MatchString(strconv.FormatFloat(val, 'f', -1, 32)) {
				return false, fmt.Errorf("%q-property: \"%f\" value doesn't match %q regular expression", tp.Key, val, restr.Restr)
			}
		}
	}
	return true, nil
}

// Parses underlying data of p as string data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsString(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertString(p)
	if !ok {
		return false, fmt.Errorf("%q-property: %q value doesn't match \"string\" data type", tp.Key, val)
	}
	for _, restr := range tp.Restrs {
		switch restr.RestrTyp {
		case TValue:
			if val != restr.Restr {
				return false, fmt.Errorf("%q-property: %q value doesn't match %q restriction", tp.Key, val, restr.Restr)
			}
		case TRegExp:
			if !restr.Restr.(*regexp.Regexp).MatchString(val) {
				return false, fmt.Errorf("%q-property: %q value doesn't match %q regular expression", tp.Key, val, restr.Restr)
			}
		}
	}
	return true, nil
}

// Parses underlying data of p as bool data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsBool(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertBool(p)
	if !ok {
		return false, fmt.Errorf("%q-property: \"%v\" value doesn't match \"bool\" data type", tp.Key, val)
	}
	for _, restr := range tp.Restrs {
		if restr.RestrTyp == TValue && val != restr.Restr {
			return false, fmt.Errorf("%q-property: \"%v\" value doesn't match \"%v\" restriction", tp.Key, val, restr.Restr)
		}
	}
	return true, nil
}

// Parses underlying data of p as datetime data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsDateTime(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertDateTime(p)
	if !ok {
		return false, fmt.Errorf("%q-property: %q value doesn't match \"datetime\" data type", tp.Key, val)
	}
	for _, restr := range tp.Restrs {
		switch restr.RestrTyp {
		case TValue:
			if val != restr.Restr {
				return false, fmt.Errorf("%q-property: %q value doesn't match %q restriction", tp.Key, val, restr.Restr)
			}
		case TRegExp:
			if !restr.Restr.(*regexp.Regexp).MatchString(val.Format(time.RFC3339)) {
				return false, fmt.Errorf("%q-property: %q value doesn't match %q regular expression", tp.Key, val, restr.Restr)
			}
		}
	}
	return true, nil
}

// Parses underlying data of p as array data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsArr(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertArr(tp.ValTyp, p)
	if !ok {
		return false, fmt.Errorf("%q-property: %q value doesn't match \"array\" data type", tp.Key, val)
	}
	for _, restr := range tp.Restrs {
		switch restr.RestrTyp {
		case TValue:
			var match bool
			rval := reflect.ValueOf(val)
			for i := 0; i < rval.Len(); i++ {
				val := rval.Index(i).Interface()
				if val == restr.Restr {
					match = true
				}
			}
			if !match {
				return false, fmt.Errorf("%q-property: array doesn't contain any value that matches %q restriction", tp.Key, restr.Restr)
			}
		case TRegExp:
			rval := reflect.ValueOf(val)
			for i := 0; i < rval.Len(); i++ {
				val := rval.Index(i).String()
				if !restr.Restr.(*regexp.Regexp).MatchString(val) {
					return false, fmt.Errorf("%q-property: %q array-value doesn't match %q regular expression", tp.Key, val, restr.Restr)
				}
			}
		}
	}
	return true, nil
}

// Parses underlying data of p as map data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsMap(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertMap(tp.KeyTyp, tp.ValTyp, p)
	if !ok {
		return false, fmt.Errorf("%q-property: %q value doesn't match \"map\" data type", tp.Key, val)
	}
	for _, restr := range tp.Restrs {
		switch restr.RestrTyp {
		case TValue:
			var match bool
			rval := reflect.ValueOf(val).MapRange()
			for rval.Next() {
				val := rval.Value().Interface()
				if val == restr.Restr {
					match = true
				}
			}
			if !match {
				return false, fmt.Errorf("%q-property: map doesn't contain any value that matches %q restriction", tp.Key, restr.Restr)
			}
		case TRegExp:
			rval := reflect.ValueOf(val).MapRange()
			for rval.Next() {
				val := rval.Value().String()
				if !restr.Restr.(*regexp.Regexp).MatchString(val) {
					return false, fmt.Errorf("%q-property: %q map-value under %q key doesn't match %q regular expression", tp.Key, val, rval.Key(), restr.Restr)
				}
			}
		case TKeyValue:
			var match bool
			rval := reflect.ValueOf(val).MapRange()
			for rval.Next() {
				val := rval.Key().Interface()
				if val == restr.Restr {
					match = true
				}
			}
			if !match {
				return false, fmt.Errorf("%q-property: map doesn't contain any key that matches %q restriction", tp.Key, restr.Restr)
			}
		case TKeyRegExp:
			rval := reflect.ValueOf(val).MapRange()
			for rval.Next() {
				val := rval.Key().String()
				if !restr.Restr.(*regexp.Regexp).MatchString(val) {
					return false, fmt.Errorf("%q-property: %q map-key doesn't match %q regular expression", tp.Key, val, restr.Restr)
				}
			}
		}
	}
	return true, nil
}

// ------------- UNKNOWN TYPE ASSERTATION ------------- //

// Searches Node- and Edge-structs within TemplateHolder-struct using v
// underlying data and returns them if found; otherwise returns nil
//
// WARNING: under the hood func tries to find types name of the v underlying
// data firstly within type definition and then within v properties values;
// thus this func can evaluate ONLY structs, maps and map-based custom types
func (t TemplateHolder) searchUnknownTypeName(v reflect.Value) (*TNode, *TEdge) {
	typName := v.Type().Name()
	asNode := t.Nodes[typName]
	asEdge := t.Edges[typName]
	if asNode == nil && asEdge == nil {
		switch v.Kind() {
		case reflect.Map:
			for iter := v.MapRange(); iter.Next(); {
				pTyp := iter.Value().String()
				if asNode == nil {
					asNode = t.Nodes[pTyp]
				}
				if asEdge == nil {
					asEdge = t.Edges[pTyp]
				}
			}
		case reflect.Struct:
			for i := 0; i < v.Type().NumField(); i++ {
				pTyp := v.Field(i).String()
				if asNode == nil {
					asNode = t.Nodes[pTyp]
				}
				if asEdge == nil {
					asEdge = t.Edges[pTyp]
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

// Parses and compares ps properties (as origin) and v underlying data (considered
// as map) properties (as implementation of origin) and merges them into a single
// "valid" map (where v keys is mapped to ps keys) of keys which is returned on
// success; omit argument is used when v underlying data contains its own type name
// as one the properties value and thus should be omitted - otherwise it can be ""
func evaluateUnknownMapPropertyKeys(omit string, v reflect.Value, ps map[string]*TProperty) (map[string]string, error) {
	vKeys := make([]string, 0)
	for iter := v.MapRange(); iter.Next(); {
		if omit == iter.Value().String() {
			continue
		}
		vKeys = append(vKeys, iter.Key().String())
	}
	pKeys := make([]string, 0, len(ps))
	for k := range ps {
		pKeys = append(pKeys, k)
	}
	validKeys, err := evaluatePropertyKeys(pKeys, vKeys)
	if err != nil {
		return nil, err
	}
	return validKeys, nil
}

// Parses and compares ps properties (as origin) and v underlying data (considered
// as struct) properties (as implementation of origin) and merges them into a single
// "valid" map (where v keys is mapped to ps keys) of keys which is returned on success;
// omit argument is used when v underlying data contains its own type name as one the
// properties value and thus should be omitted - otherwise it can be ""
func evaluateUnknownStructPropertyKeys(omit string, v reflect.Value, ps map[string]*TProperty) (map[string]string, error) {
	vKeys := make([]string, 0)
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		if omit == v.Field(i).String() {
			continue
		}
		vKeys = append(vKeys, vt.Field(i).Name)
	}
	pKeys := make([]string, 0, len(ps))
	for k := range ps {
		pKeys = append(pKeys, k)
	}

	validKeys, err := evaluatePropertyKeys(pKeys, vKeys)
	if err != nil {
		return nil, err
	}
	return validKeys, nil
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

// -------------- BASE TYPES ASSERTATION -------------- //

// Asserts - is the data type of the underlying value of v is int; returns
// asserted int value and true on success
func assertInt(v interface{}) (int, bool) {
	val, ok := v.(int)
	return val, ok
}

// Asserts - is the data type of the underlying value of v is float; returns
// asserted float value and true on success
func assertFloat(v interface{}) (float64, bool) {
	switch res := v.(type) {
	case float64:
		return res, true
	case float32:
		return float64(res), true
	default:
		return 0.0, false
	}
}

// Asserts - is the data type of the underlying value of v is string; returns
// asserted string value and true on success
func assertString(v interface{}) (string, bool) {
	val, ok := v.(string)
	return val, ok
}

// Asserts - is the data type of the underlying value of v is bool; returns
// asserted bool value and true on success
func assertBool(v interface{}) (bool, bool) {
	val, ok := v.(bool)
	return val, ok
}

// Asserts - is the data type of the underlying value of v is datetime; returns
// asserted datetime value and true on success
func assertDateTime(v interface{}) (time.Time, bool) {
	val, ok := v.(time.Time)
	return val, ok
}

// Asserts - is the data type of the underlying value of v is array with vt
// data type of "inner" values; returns asserted array value wrapped into blank
// interface and true on success
func assertArr(vt TDataType, p interface{}) (interface{}, bool) {
	i := reflect.ValueOf(p)
	rvt := i.Index(0)
	switch vt {
	case TInt:
		if rvt.Kind() != reflect.Int {
			return nil, false
		}
	case TFloat:
		if rvt.Kind() != reflect.Float64 {
			return nil, false
		}
	case TString:
		if rvt.Kind() != reflect.String {
			return nil, false
		}
	case TBool:
		if rvt.Kind() != reflect.Bool {
			return nil, false
		}
	case TDateTime:
		if _, ok := rvt.Interface().(time.Time); !ok {
			return nil, false
		}
	}
	return p, true
}

// Asserts - is the data type of the underlying value of v is map with vt
// data type of "inner" values and kt data type of value keys; returns
// asserted map value wrapped into blank interface and true on success
func assertMap(kt, vt TDataType, p interface{}) (interface{}, bool) {
	i := reflect.ValueOf(p).MapRange()
	i.Next()
	rkt := i.Key()
	rvt := i.Value()
	switch kt {
	case TInt:
		if rkt.Kind() != reflect.Int {
			return nil, false
		}
	case TFloat:
		if rkt.Kind() != reflect.Float64 {
			return nil, false
		}
	case TString:
		if rkt.Kind() != reflect.String {
			return nil, false
		}
	case TBool:
		if rkt.Kind() != reflect.Bool {
			return nil, false
		}
	case TDateTime:
		if _, ok := rkt.Interface().(time.Time); !ok {
			return nil, false
		}
	}
	switch vt {
	case TInt:
		if rvt.Kind() != reflect.Int {
			return nil, false
		}
	case TFloat:
		if rvt.Kind() != reflect.Float64 {
			return nil, false
		}
	case TString:
		if rvt.Kind() != reflect.String {
			return nil, false
		}
	case TBool:
		if rvt.Kind() != reflect.Bool {
			return nil, false
		}
	case TDateTime:
		if _, ok := rvt.Interface().(time.Time); !ok {
			return nil, false
		}
	}
	return p, true
}
