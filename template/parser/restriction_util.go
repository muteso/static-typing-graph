package parser

import (
	"fmt"
	"regexp"
	"stg/template"
	"strconv"
	"strings"
	"time"
)

// regexps for validation of unprocessed (in string form) data types
var (
	intRe      = regexp.MustCompile(`^\d+$`).MatchString
	floatRe    = regexp.MustCompile(`^\d+\.\d+$`).MatchString
	datetimeRe = regexp.MustCompile(`^\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ$`).MatchString
	// string and bool data types don't need regexp
)

// Buffer type for representation of each unique data type
// combination (mostrly for maps and arrays)
type typeBuffer struct {
	T  template.TDataType
	Vt template.TDataType
	Kt template.TDataType
}

// Buffer type which stores funcs for proper mutation of any
// "simple" (int, float, string, bool, datetime) data type
type mutationTool struct {
	check  func(string) bool
	mutate func(string) (interface{}, error)
}

// List of baked mutationTool-structs for "simple" data types
var (
	intTool = mutationTool{
		check: func(v string) bool { return intRe(v) },
		mutate: func(v string) (interface{}, error) {
			i, _ := strconv.Atoi(v)
			// ignores error cus check-func validates that v is int
			return i, nil
		},
	}
	floatTool = mutationTool{
		check: func(v string) bool { return floatRe(v) },
		mutate: func(v string) (interface{}, error) {
			f, _ := strconv.ParseFloat(v, 64)
			// ignores error cus check-func validates that v is float
			return f, nil
		},
	}
	stringTool = mutationTool{
		check: func(v string) bool { return true },
		mutate: func(v string) (interface{}, error) {
			return v, nil // v already is string
		},
	}
	boolTool = mutationTool{
		check: func(v string) bool { return v == "true" || v == "false" },
		mutate: func(v string) (interface{}, error) {
			b, _ := strconv.ParseBool(v)
			// ignores error cus check-func validates that v is bool
			return b, nil
		},
	}
	datetimeTool = mutationTool{
		check: func(v string) bool { return datetimeRe(v) },
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
func getMutationTool(t string, rt template.TRestrictionType) mutationTool {
	typs := strings.Split(t, "-")
	typ := typs[0]

	switch typ {
	case "int":
		return intTool
	case "float":
		return floatTool
	case "string":
		return stringTool
	case "bool":
		return boolTool
	case "datetime":
		return datetimeTool
	case "array":
		return getMutationTool(typs[1], 0)
	case "map":
		var searchTyp string
		switch rt {
		case template.TValue:
			searchTyp = typs[2]
		case template.TRegExp:
			searchTyp = typs[2]
		case template.TKeyValue:
			searchTyp = typs[1]
		case template.TKeyRegExp:
			searchTyp = typs[1]
		}
		return getMutationTool(searchTyp, 0)
	}
	return mutationTool{}
}

// Returns according DataType-const to unprocessed (in string
// form) t data type; returns Null-const if t is incorrect or
// not a "simple" (int, float, string, bool, datetime) data type
func toSimpleDataType(t string) template.TDataType {
	switch t {
	case "int":
		return template.TInt
	case "float":
		return template.TFloat
	case "string":
		return template.TString
	case "bool":
		return template.TBool
	case "datetime":
		return template.TDateTime
	}
	return template.TNull
}

// Returns typeBuffer-struct which stores full representation of
// unprocessed (in string form) t data type; returns empty struct
// and error as reustl if t is incorrect
func toDataType(t string) (typeBuffer, error) {
	typs := strings.Split(t, "-")
	typ := typs[0]

	switch typ {
	case "int":
		if len(typs) > 1 {
			return typeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return typeBuffer{
			T:  template.TInt,
			Vt: template.TInt,
		}, nil
	case "float":
		if len(typs) > 1 {
			return typeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return typeBuffer{
			T:  template.TFloat,
			Vt: template.TFloat,
		}, nil
	case "string":
		if len(typs) > 1 {
			return typeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return typeBuffer{
			T:  template.TString,
			Vt: template.TString,
		}, nil
	case "bool":
		if len(typs) > 1 {
			return typeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return typeBuffer{
			T:  template.TBool,
			Vt: template.TBool,
		}, nil
	case "datetime":
		if len(typs) > 1 {
			return typeBuffer{}, fmt.Errorf("data type %q can't has subtypes", typ)
		}
		return typeBuffer{
			T:  template.TDateTime,
			Vt: template.TDateTime,
		}, nil
	case "array":
		if len(typs) > 2 {
			return typeBuffer{}, fmt.Errorf("data type \"array\" can't has more than 1 subtype")
		}
		vty := toSimpleDataType(typs[1])
		if vty == template.TNull {
			return typeBuffer{}, fmt.Errorf("data type \"array\" has wrong value data subtype %q", typs[1])
		}
		return typeBuffer{
			T:  template.TArray,
			Vt: vty,
		}, nil
	case "map":
		if len(typs) > 3 {
			return typeBuffer{}, fmt.Errorf("data type \"map\" can't has more than 2 subtypes - 1 for keys and 1 for values")
		}
		kty, vty := toSimpleDataType(typs[1]), toSimpleDataType(typs[2])
		eText := ""
		if kty == template.TNull {
			eText += fmt.Sprintf("data type \"map\" has wrong key data subtype %q", typs[1])
		}
		if vty == template.TNull {
			if kty == template.TNull {
				eText += "; "
			}
			eText += fmt.Sprintf("\"map\" data type has wrong value data subtype %q", typs[2])
		}
		if eText != "" {
			return typeBuffer{}, fmt.Errorf(eText)
		}
		return typeBuffer{
			T:  template.TMap,
			Kt: kty,
			Vt: vty,
		}, nil
	}
	return typeBuffer{}, fmt.Errorf("undefined data type %q", t)
}

// Mutates r to actual temaplate Restriction-struct using t and rt
// to correct mutation. If any error occurs or t and rt conflicts
// with each other - returns nil and error as result
func mutateRestr(t typeBuffer, rt template.TRestrictionType, r string) (*template.TRestriction, error) {
	if t.T == template.TNull {
		return nil, fmt.Errorf("restriction %q can't be inferred because of undefined or wrong data type of restricted property", r)
	}
	actual := &template.TRestriction{
		Typ:      t.Vt,
		RestrTyp: rt,
	}

	if rt == template.TKeyValue || rt == template.TKeyRegExp {
		if t.T != template.TMap {
			return nil, fmt.Errorf("data type %q can't has key restrictions", t.T)
		}
		actual.Typ = t.Kt
	}

	if rt == template.TRegExp || rt == template.TKeyRegExp {
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
func checkRestrContradiction(re *template.TRestriction, rs []*template.TRestriction) (bool, error) {
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
		case template.TInt:
			str = strconv.FormatInt(int64(v.(int)), 10)
		case template.TFloat:
			str = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		case template.TString:
			str = v.(string)
		case template.TBool:
			str = strconv.FormatBool(v.(bool))
		case template.TDateTime:
			str = v.(time.Time).Format(time.RFC3339)
		}
		if ok := regexp.MatchString(str); !ok {
			return ok, fmt.Errorf("%q regexp contradicts %q value restriction", regexp.String(), str)
		}
	}
	return true, nil
}
