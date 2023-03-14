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
// and returns nil if doesn't match any contradictions between them
func comparePropertyKeys(tKeys, vKeys []string) error {
	tempKeys := make(map[string]string)
	for _, k := range tKeys {
		tempKeys[k] = k
	}
	for _, k := range vKeys {
		if _, ok := tempKeys[k]; !ok {
			return fmt.Errorf("validated entity has extra %q property", k)
		} else {
			delete(tempKeys, k)
		}
	}
	if len(tempKeys) != 0 {
		missedProps := ""
		for _, k := range tempKeys {
			missedProps += fmt.Sprintf("%q, ", k)
		}
		missedProps = strings.TrimSuffix(missedProps, ", ")
		return fmt.Errorf("validated entity doesn't has %s properties", missedProps)
	}
	return nil
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
		return false, fmt.Errorf("%q-property: \"%d\" value doesn't match \"int\" data type", tp.Key, val)
	}
	if len(tp.ValRestrs) == 0 {
		return true, nil
	}
	for _, restr := range tp.ValRestrs {
		switch restr.RestrTyp {
		case TValue:
			if val == restr.Restr {
				return true, nil
			}
		case TRegExp:
			if restr.Restr.(*regexp.Regexp).MatchString(strconv.FormatInt(int64(val), 10)) {
				return true, nil
			}
		}
	}
	return false, fmt.Errorf("%q-property: \"%d\" value doesn't match neither restrictions", tp.Key, val)
}

// Parses underlying data of p as float data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsFloat(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertFloat(p)
	if !ok {
		return false, fmt.Errorf("%q-property: value \"%f\" doesn't match \"float\" data type", tp.Key, val)
	}
	if len(tp.ValRestrs) == 0 {
		return true, nil
	}
	for _, restr := range tp.ValRestrs {
		switch restr.RestrTyp {
		case TValue:
			if val == restr.Restr {
				return true, nil
			}
		case TRegExp:
			if restr.Restr.(*regexp.Regexp).MatchString(strconv.FormatFloat(val, 'f', -1, 32)) {
				return true, nil
			}
		}
	}
	return false, fmt.Errorf("%q-property: \"%f\" value doesn't match neither restrictions", tp.Key, val)
}

// Parses underlying data of p as string data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsString(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertString(p)
	if !ok {
		return false, fmt.Errorf("%q-property: %q value doesn't match \"string\" data type", tp.Key, val)
	}
	if len(tp.ValRestrs) == 0 {
		return true, nil
	}
	for _, restr := range tp.ValRestrs {
		switch restr.RestrTyp {
		case TValue:
			if val == restr.Restr {
				return true, nil
			}
		case TRegExp:
			if restr.Restr.(*regexp.Regexp).MatchString(val) {
				return true, nil
			}
		}
	}
	return false, fmt.Errorf("%q-property: %q value doesn't match neither restrictions", tp.Key, val)
}

// Parses underlying data of p as bool data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsBool(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertBool(p)
	if !ok {
		return false, fmt.Errorf("%q-property: \"%v\" value doesn't match \"bool\" data type", tp.Key, val)
	}
	if len(tp.ValRestrs) == 0 {
		return true, nil
	}
	for _, restr := range tp.ValRestrs {
		if restr.RestrTyp == TValue && val == restr.Restr {
			return true, nil
		}
	}
	return false, fmt.Errorf("%q-property: \"%v\" value doesn't match neither restrictions", tp.Key, val)
}

// Parses underlying data of p as datetime data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsDateTime(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertDateTime(p)
	if !ok {
		return false, fmt.Errorf("%q-property: %q value doesn't match \"datetime\" data type", tp.Key, val)
	}
	if len(tp.ValRestrs) == 0 {
		return true, nil
	}
	for _, restr := range tp.ValRestrs {
		switch restr.RestrTyp {
		case TValue:
			if val == restr.Restr {
				return true, nil
			}
		case TRegExp:
			if restr.Restr.(*regexp.Regexp).MatchString(val.Format(time.RFC3339)) {
				return true, nil
			}
		}
	}
	return false, fmt.Errorf("%q-property: %q value doesn't match neither restrictions", tp.Key, val)
}

// Parses underlying data of p as array data type property and validates it using
// tp as validator and returns true and nil on success
func evaluatePropertyAsArr(tp TProperty, p interface{}) (bool, error) {
	val, ok := assertArr(tp.ValTyp, p)
	if !ok {
		return false, fmt.Errorf("%q-property: %q value doesn't match \"array\" data type", tp.Key, val)
	}
	if len(tp.ValRestrs) == 0 {
		return true, nil
	}
	for i := 0; i < val.Len(); i++ {
		val := val.Index(i)
		match := false
		for _, restr := range tp.ValRestrs {
			switch restr.RestrTyp {
			case TValue:
				if val.Interface() == restr.Restr {
					match = true
					break
				}
			case TRegExp:
				if restr.Restr.(*regexp.Regexp).MatchString(val.String()) {
					match = true
					break
				}
			}
		}
		if !match {
			return false, fmt.Errorf("%q-property: %q value doesn't match neither of the restrictions", tp.Key, val.String())
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
	if len(tp.ValRestrs) == 0 && len(tp.KeyRestrs) == 0 {
		return true, nil
	}
	// matchCountValues, matchCountKeys := 0, 0
	for iter := val.MapRange(); iter.Next(); {
		val := iter.Value()
		matchValue := false
		for _, restr := range tp.ValRestrs {
			switch restr.RestrTyp {
			case TValue:
				if val.Interface() == restr.Restr {
					matchValue = true
					break
				}
			case TRegExp:
				if restr.Restr.(*regexp.Regexp).MatchString(val.String()) {
					matchValue = true
					break
				}
			}
		}
		if !matchValue {
			return false, fmt.Errorf("%q-property: %q value doesn't match neither of the restrictions", tp.Key, val.String())
		}
		key := iter.Key()
		matchKey := false
		for _, restr := range tp.KeyRestrs {
			switch restr.RestrTyp {
			case TKeyValue:
				if key.Interface() == restr.Restr {
					matchKey = true
					break
				}
			case TKeyRegExp:
				if restr.Restr.(*regexp.Regexp).MatchString(key.String()) {
					matchKey = true
					break
				}
			}
		}
		if !matchKey {
			return false, fmt.Errorf("%q-property: %q kay doesn't match neither of the restrictions", tp.Key, key.String())
		}
	}

	// if matchCountValues < val.Len() {
	// 	return false, fmt.Errorf("%q-property: not all map values satisfy restrictions", tp.Key)
	// }
	// if matchCountKeys < val.Len() {
	// 	return false, fmt.Errorf("%q-property: not all map keys satisfy restrictions", tp.Key)
	// }
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
// data type of "inner" values; returns asserted array value wrapped into
// reflect.Value interface and true on success
func assertArr(vt TDataType, v interface{}) (reflect.Value, bool) {
	i := reflect.ValueOf(v)
	if i.Kind() != reflect.Slice && i.Kind() != reflect.Array {
		return reflect.Value{}, false
	}
	rv := i.Index(0)
	if !assertSimpleDataType(vt, rv) {
		return reflect.Value{}, false
	}
	return i, true
}

// Asserts - is the data type of the underlying value of v is map with vt
// data type of "inner" values and kt data type of value keys; returns
// asserted map value wrapped into blank reflect.Value and true on success
func assertMap(kt, vt TDataType, v interface{}) (reflect.Value, bool) {
	i := reflect.ValueOf(v)
	if i.Kind() != reflect.Map {
		return reflect.Value{}, false
	}
	iter := i.MapRange()
	iter.Next()
	rk := iter.Key()
	rv := iter.Value()
	if !assertSimpleDataType(kt, rk) {
		return reflect.Value{}, false
	}
	if !assertSimpleDataType(vt, rv) {
		return reflect.Value{}, false
	}
	return i, true
}

// Asserts - is the data type of the underlying value of v is a "simple"
// data type; returns true on success
func assertSimpleDataType(t TDataType, v reflect.Value) bool {
	prt := v.Kind()
	switch t {
	case TInt:
		return prt == reflect.Int
	case TFloat:
		return prt == reflect.Float64
	case TString:
		return prt == reflect.String
	case TBool:
		return prt == reflect.Bool
	case TDateTime:
		_, ok := v.Interface().(time.Time)
		return ok
	}
	return true
}
