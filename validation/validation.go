package validation

// Creates and returns new struct that implements Node- and
// Edge-interfaces with typ type name and props properties
//
// This func conveniently wraps into NewNode- or NewEdge-func
func newEntity(typ string, props map[string]interface{}) entityHolder {
	return entityHolder{
		typ:   typ,
		props: props,
	}
}

// Creates and returns new Node-interface value with typ type name
// and props properties
func NewNode(typ string, props map[string]interface{}) Node {
	return newEntity(typ, props)
}

// Creates and returns new Edge-interface value with typ type name
// and props properties
func NewEdge(typ string, props map[string]interface{}) Edge {
	return newEntity(typ, props)
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

// Validates v (which might implements Node-, Edge- or Triplet-interface)
// underlying data using vr and returns true and nil on success; in cases
// where underlying data dont implement interfaces enumerated above this
// func still tries to validate it but in some successive cases may also
// return non-nil error explaning what was exactly validated
func Validate(vr Validator, v interface{}) (bool, error) {
	switch val := v.(type) {
	case Triplet:
		return vr.ValidateTriplet(val)
	case Node:
		return vr.ValidateNode(val)
	case Edge:
		return vr.ValidateEdge(val)
	default:
		return vr.ValidateUnknown(val)
	}
}
