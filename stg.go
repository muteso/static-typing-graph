/*
 */
package stg

import (
	"io"
	"stg/template/parser"
	"stg/validation"
)

// Type shortcuts
type (
	// Validator interface that is used as validation tool
	Validator = validation.Validator
	// Node interface - represents graph node which can return his type's
	// name and properties
	Node = validation.Node
	// Edge interface - represents graph edge which can return his type's
	// name and properties
	Edge = validation.Edge
	// Triplet interface that uses RDF-representation as sematic reference;
	// any type that can return main entity, subject entity and entity, that
	// semanticly connects them, implements that interface
	Triplet = validation.Triplet
)

// Parses ALREADY opened template-file (or any another representation
// of it implementing io.Reader-interface) and returns result
// as Validator-interface value; if any error occurs doesn't interrupt
// parsing and then returns error which contains all occured errors during
// parsing
func ParseTemplate(file io.Reader) (Validator, error) {
	return parser.ParseTemplate(file)
}

// Creates and returns new Node-interface value with typ type name
// and props properties
func NewNode(typ string, props map[string]interface{}) Node {
	return validation.NewNode(typ, props)
}

// Creates and returns new Edge-interface value with typ type name
// and props properties
func NewEdge(typ string, props map[string]interface{}) Edge {
	return validation.NewEdge(typ, props)
}

// Creates and returns new Triplet-interface with m as main node entity,
// s as subject node entity and e as edge entity
func NewTriplet(m, s Node, e Edge) Triplet {
	return validation.NewTriplet(m, s, e)
}

// Validates v (which might implements Node-, Edge- or Triplet-interface)
// underlying data using vr and returns true and nil on success; in cases
// where underlying data dont implement interfaces enumerated above this
// func still tries to validate it but in some successive cases may also
// return non-nil error explaning what was exactly validated
func Validate(vr Validator, v interface{}) (bool, error) {
	return validation.Validate(vr, v)
}
