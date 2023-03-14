/*
stg-package (Static Typing for Graph data models) provides convinient instruments
for wrapping almost any possible go-data into graph-data and VALIDATING it, as it
has static type (with some additional constraints) and being part of a graph! For
validating itself you should first define template (which, to be honest, is the most
difficult part of the work) in the yaml-file which then will be parsed and used as
validator-tool.
But there are some limitations:
  - maps can't nest within each other (which should be handled by making
    a new node/edge that contains nested map etc.)
  - arrays can't nest within each other (for the same reason as above)
  - maps can't nest within arrays and vise versa (for the same reason as above)
  - map's and array's values (and keys in case of maps) can contain only the
    same PRIMITIVE (int, string, etc) data types (for reasons of go data
    types compatibility)
  - labels can't nest within each other (which may lead to endless data-nesting and
    implicit property- and connection-definitions which may cause hard-to-find
    definition conflicts)
  - graph-data can't contain 2 or more identical nodes (only one instance of node
    will be created)
  - graph-data can't contain 2 or more identical edges with identical directions
    between one unique pair of nodes (only one instance of edge will be created)

Otherwise, public interface contains only few types and functions which should make
ease to use this library!

Types:

	Validator
	Node
	Edge
	Triplet
	Duplet
	Graph

Functions:

	ParseTemplate(template-file) Validator
	NewNode(type, properties) Node
	NewEdge(type, properties) Edge
	NewTriplet(main node, subject node, edge) Triplet
	NewDuplet(node, edge) Duplet
	NewGraph(triplets) Graph
	Validate(validator, any graph entity) bool, error

As simple as it looks!
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
	// Triplet interface that uses RDF-representation as semantic reference;
	// any type that can return main entity, subject entity and entity, that
	// semanticly connects them, implements that interface
	Triplet = validation.Triplet
	// Duplet interface that uses "halfed" RDF-representation as semantic
	// reference - Duplet-interface express only "tail" or "head" of Triplet-interface
	// (without main entity or without subject entity accordingly); any type
	// that can return "head"/"tail" entity and entity, that semanticly connects
	// "head"/"tail" entity with "invisible" "tail"/"head" entity, implements
	// that interface
	Duplet = validation.Duplet
	// Graph interface - represents some amount of nodes interconnected (or
	// not) with ecah other by edges; have methods for both introspection and
	// mutation
	//
	// WARNING: Graph-interface should NOT assign some sort of ID's to the nodes
	// and edges and should NOT do some sort of indexing, so it can't be used as
	// fully functional inmemory DB; also, because of this Graph-interface should
	// NOT allow to contain identical nodes and edges within implementing data type
	Graph = validation.Graph
)

// Parses ALREADY opened template-file (or any another representation
// of it implementing io.Reader-interface) and returns result
// as Validator-interface value; if any error occurs doesn't interrupt
// parsing and then returns error which contains all occured errors during
// parsing
//
// WARNING: template dont allow some graph design practices, as:
//   - maps nesting (which should be handled by making a new node/edge that
//     contains nested map etc.)
//   - arrays nesting (for the same reason as above)
//   - nesting maps within arrays and vise versa (for the same reason as above)
//   - map's and array's values (and keys in case of maps) can contain only the
//     same PRIMITIVE (int, string, etc) data types (for reasons of go data
//     types compatibility)
//   - labels can't nest within each other (which may lead to endless data-nesting
//     and implicit property- and connection-definitions which may cause hard-to-find
//     definition conflicts)
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

// Creates and returns new Duplet-interface with n as "tail"/"head" node
// entity and e as edge entity
func NewDuplet(n Node, e Edge) Duplet {
	return validation.NewDuplet(n, e)
}

// Creates and returns new Graph-interface using ns nodes and (optionally,
// but more recommended) gr triplets for building it - this func parses any
// given node and triplet and inserts it into Graph-interface, interconnecting
// "newcoming" entities with already inserted graph-entities (using main nodes
// of triplets for infering where it should happen in case of triplets)
//
// WARNING: this func dont allow some graph design practices that leads
// to excess data storaging, as:
//   - 2 or more identical nodes (this func will create only one instance of node)
//   - 2 or more identical edges with identical directions between one unique pair
//     of nodes (this func will create only one instance of edge)
//   - identity of nodes (and edges too) is found if they have equal types and
//     ALL their properties are equal
//
// Respectively, triplets with equal main nodes will be "squashed" into a
// "single" graph node and any duplicates will be just OMITTED (i.e. identical
// triplets, nodes and edges)
//
// Also, in some cases inserting of triplets has a little bit different semanthic:
//   - already presented edge - cancels insertation
//   - already presented main node and subject node - connects them with
//     edge (adding edge)
//   - nil subject node and nil edge - inserts main node as single node
//     (adding node)
//   - nil main node and nil edge - inserts subject node as single node
//     (adding node)
//   - nil edge - inserts main node and subject node as single nodes (adding
//     nodes)
//
// Because of this it's more recommendend to use ONLY triplets to build graphs, which
// may be done by using:
//
//	graph := NewGraph(nil, tripletSlice...)
func NewGraph(ns []Node, gr ...Triplet) Graph {
	return validation.NewGraph(ns, gr...)
}

// Validates v (which might implements Node-, Edge- or Triplet-interface)
// underlying data using vr and returns true and nil on success; in cases
// where underlying data dont implement interfaces enumerated above this
// func still tries to validate it but in some successive cases may also
// return non-nil error explaning what was exactly validated
func Validate(vr Validator, v interface{}) (bool, error) {
	return validation.Validate(vr, v)
}
