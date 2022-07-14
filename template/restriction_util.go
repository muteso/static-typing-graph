package template

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// regexps for validation of unprocessed (in string form) data types
var (
	int_re      = regexp.MustCompile(`^\d+$`).MatchString
	float_re    = regexp.MustCompile(`^\d+\.\d+$`).MatchString
	datetime_re = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`).MatchString
	// string and bool data types don't need regexp
)

// Buffer type for representation of each unique data type
// combination (mostrly for maps and arrays)
type TypeBuffer struct {
	T  TDataType
	Vt TDataType
	Kt TDataType
}

// Buffer type which stores funcs for proper mutation of any
// "simple" (int, float, string, bool, datetime) data type
type mutationTool struct {
	check  func(string) bool
	mutate func(string) (interface{}, error)
}

// List of baked mutationTool-structs for "simple" data types
var (
	int_tool = mutationTool{
		check: func(v string) bool { return int_re(v) },
		mutate: func(v string) (interface{}, error) {
			i, _ := strconv.Atoi(v)
			// ignores error cus check-func validates that v is int
			return i, nil
		},
	}
	float_tool = mutationTool{
		check: func(v string) bool { return float_re(v) },
		mutate: func(v string) (interface{}, error) {
			f, _ := strconv.ParseFloat(v, 64)
			// ignores error cus check-func validates that v is float
			return f, nil
		},
	}
	string_tool = mutationTool{
		check: func(v string) bool { return true },
		mutate: func(v string) (interface{}, error) {
			return v, nil // v already is string
		},
	}
	bool_tool = mutationTool{
		check: func(v string) bool { return v == "true" || v == "false" },
		mutate: func(v string) (interface{}, error) {
			b, _ := strconv.ParseBool(v)
			// ignores error cus check-func validates that v is bool
			return b, nil
		},
	}
	datetime_tool = mutationTool{
		check: func(v string) bool { return datetime_re(v) },
		mutate: func(v string) (interface{}, error) {
			d, e := time.Parse(time.RFC3339, v)
			if e != nil {
				return nil, fmt.Errorf("\"datetime\" restriction doesn't match RFC3339 (YYYY-MM-DDTHH:MM:SSZ) format: %s", e)
			}
			return d, nil
		},
	}
)

// Returns mutationTool-struct accordingly to unprocessed (in
// string form) t data type; returns empty struct if t is incorrect
func getMutationTool(t string, rt TRestrictionType) mutationTool {
	typs := strings.Split(t, "-")
	typ := typs[0]

	switch typ {
	case "int":
		return int_tool
	case "float":
		return float_tool
	case "string":
		return string_tool
	case "bool":
		return bool_tool
	case "datetime":
		return datetime_tool
	case "array":
		return getMutationTool(typs[1], 0)
	case "map":
		var search_typ string
		switch rt {
		case TValue:
			search_typ = typs[2]
		case TRegExp:
			search_typ = typs[2]
		case TKeyValue:
			search_typ = typs[1]
		case TKeyRegExp:
			search_typ = typs[1]
		}
		return getMutationTool(search_typ, 0)
	}
	return mutationTool{}
}

// Returns according DataType-const to unprocessed (in string
// form) t data type; returns Null-const if t is incorrect or
// not a "simple" (int, float, string, bool, datetime) data type
func toSimpleDataType(t string) TDataType {
	switch t {
	case "int":
		return TInt
	case "float":
		return TFloat
	case "string":
		return TString
	case "bool":
		return TBool
	case "datetime":
		return TDateTime
	}
	return TNull
}

// Returns TypeBuffer-struct which stores full representation of
// unprocessed (in string form) t data type; returns empty struct
// and error as reustl if t is incorrect
func ToDataType(t string) (TypeBuffer, error) {
	typs := strings.Split(t, "-")
	typ := typs[0]

	switch typ {
	case "int":
		if len(typs) > 1 {
			return TypeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return TypeBuffer{
			T:  TInt,
			Vt: TInt,
		}, nil
	case "float":
		if len(typs) > 1 {
			return TypeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return TypeBuffer{
			T:  TFloat,
			Vt: TFloat,
		}, nil
	case "string":
		if len(typs) > 1 {
			return TypeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return TypeBuffer{
			T:  TString,
			Vt: TString,
		}, nil
	case "bool":
		if len(typs) > 1 {
			return TypeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return TypeBuffer{
			T:  TBool,
			Vt: TBool,
		}, nil
	case "datetime":
		if len(typs) > 1 {
			return TypeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return TypeBuffer{
			T:  TDateTime,
			Vt: TDateTime,
		}, nil
	case "array":
		if len(typs) > 2 {
			return TypeBuffer{}, fmt.Errorf("data type \"array\" can't has more than 1 subtype")
		}
		vty := toSimpleDataType(typs[1])
		if vty == TNull {
			return TypeBuffer{}, fmt.Errorf("data type \"array\" has wrong value data subtype %q", typs[1])
		}
		return TypeBuffer{
			T:  TArray,
			Vt: vty,
		}, nil
	case "map":
		if len(typs) > 3 {
			return TypeBuffer{}, fmt.Errorf("data type \"map\" can't has more than 2 subtypes - 1 for keys and 1 for values")
		}
		kty, vty := toSimpleDataType(typs[1]), toSimpleDataType(typs[2])
		e_text := ""
		if kty == TNull {
			e_text += fmt.Sprintf("data type \"map\" has wrong key data subtype %q", typs[1])
		}
		if vty == TNull {
			if kty == TNull {
				e_text += "; "
			}
			e_text += fmt.Sprintf("\"map\" data type has wrong value data subtype %q", typs[2])
		}
		if e_text != "" {
			return TypeBuffer{}, fmt.Errorf(e_text)
		}
		return TypeBuffer{
			T:  TMap,
			Kt: kty,
			Vt: vty,
		}, nil
	}
	return TypeBuffer{}, fmt.Errorf("undefined data type %q", t)
}

// Mutates r to actual temaplate Restriction-struct using t and rt
// to correct mutation. If any error occurs or t and rt conflicts
// with each other - returns nil and error as result
func MutateRestr(t TypeBuffer, rt TRestrictionType, r string) (*TRestriction, error) {
	if t.T == TNull {
		return nil, fmt.Errorf("restriction %q can't be inferred because of undefined or wrong data type of restricted property", r)
	}
	actual := &TRestriction{
		Typ:      t.Vt,
		RestrTyp: rt,
	}

	if rt == TKeyValue || rt == TKeyRegExp {
		if t.T != TMap {
			return nil, fmt.Errorf("data type %q can't has key restrictions", t.T)
		}
		actual.Typ = t.Kt
	}

	if rt == TRegExp || rt == TKeyRegExp {
		re, err := regexp.Compile(r)
		if err != nil {
			return nil, fmt.Errorf("restriction %q has wrong regexp schema", r)
		}
		actual.Restr = re
		return actual, nil
	}

	mut := getMutationTool(actual.Typ.String(), actual.RestrTyp)
	if !mut.check(r) {
		return nil, fmt.Errorf("restriction %q doesn't match %q data type", r, t.Vt)
	}
	res, err := mut.mutate(r)
	if err != nil {
		return nil, err
	}
	actual.Restr = res
	return actual, nil
}

// Checks if the given re regexp-restriction contradicts any of the
// rs value restrictions; returns false and error if contradiction
// occurs
func CheckRestrContradiction(re *TRestriction, rs []*TRestriction) (bool, error) {
	regexp, ok := re.Restr.(*regexp.Regexp)
	if !ok {
		return false, nil
	}
	for _, restr := range rs {
		if restr == nil {
			continue
		}
		var str string
		switch v := restr.Restr; restr.Typ {
		case TInt:
			str = strconv.FormatInt(int64(v.(int)), 10)
		case TFloat:
			str = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		case TString:
			str = v.(string)
		case TBool:
			str = strconv.FormatBool(v.(bool))
		case TDateTime:
			str = v.(time.Time).Format(time.RFC3339)
		}
		if ok := regexp.MatchString(str); !ok {
			return ok, fmt.Errorf("%q regexp contradicts %q value restriction", regexp.String(), str)
		}
	}
	return true, nil
}
