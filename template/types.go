package template

const INF = -1

// Type that represents data type of template property
type TDataType uint8

const (
	TNull TDataType = iota
	TInt
	TFloat
	TString
	TBool
	TDateTime
	TArray
	TMap
)

// Common String-method to implement Stringer-interface
func (d TDataType) String() string {
	switch d {
	case TInt:
		return "int"
	case TFloat:
		return "float"
	case TString:
		return "string"
	case TBool:
		return "bool"
	case TMap:
		return "map"
	case TArray:
		return "array"
	case TDateTime:
		return "datetime"
	}
	return "null"
}

// Type that represents HOW template restricton restricts values
// of any given data type
type TRestrictionType uint8

const (
	_ TRestrictionType = iota
	TValue
	TRegExp
	TKeyValue
	TKeyRegExp
)

// Common String-method to implement Stringer-interface
func (r TRestrictionType) String() string {
	switch r {
	case TValue:
		return "value"
	case TRegExp:
		return "regexp"
	case TKeyValue:
		return "key value"
	case TKeyRegExp:
		return "key regexp"
	}
	return ""
}

// Template node type - contains type name and properties
type TNode struct {
	Typ   string
	Props map[string]*TProperty
}

// Template edge type - contains type name and properties
type TEdge struct {
	Typ   string
	Props map[string]*TProperty
}

// Template property type - represents key:value-pair; contains
// key and value's data type and restrictions; if type is "simple"
// (int, float, string, bool, datetime) this type is also counted
// as "inner" value type; if Type is Array it also contains type of
// "inner" values; if Type is Map it also contains type of map keys
type TProperty struct {
	Key       string
	Typ       TDataType
	ValTyp    TDataType
	KeyTyp    TDataType
	ValRestrs []*TRestriction
	KeyRestrs []*TRestriction
}

// Template restriction type - represents restriction of property;
// contains type of restriction and restriction itself
type TRestriction struct {
	Typ      TDataType
	RestrTyp TRestrictionType
	Restr    interface{}
}

// Template connection type - represents bound between main node and
// subject node, connected by Edge, which ALWAYS directed from main node
// to subject; contains main node, edge, subject node and minimum and
// maximum possible amount of connections between ONE unique main node
// and ANY amount of unique subject nodes using edge
type TConnection struct {
	Main *TNode
	Edge *TEdge
	Subj *TNode

	Min int
	Max int
}
