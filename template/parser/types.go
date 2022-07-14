package parser

import (
	"strings"
)

// Custom type that represents parse error (location where
// error had occured and description of what happened)
type parseError struct {
	Loc string
	Msg string
}

// Common Error-method to implement error-interface
func (e parseError) Error() string {
	return e.Loc + " >> " + e.Msg
}

// Type for embedding to structs that contains information about
// nesting within "parent" struct
type nesting []string

// Returns pretty formatted sequence of nesting in string format
func (n nesting) String() string {
	return "template | " + strings.Join(n, " | ")
}

// Temporal buffer type for .yaml parsing purposes; represents
// template-file itself
type bTemplate struct {
	BufLabels map[string]bLabel `yaml:"labels"`
	BufNodes  map[string]bNode  `yaml:"nodes"`
	BufEdges  map[string]bEdge  `yaml:"edges"`
}

// Temporal buffer type for .yaml parsing purposes; represents
// sets-field of template-file
type bLabel struct {
	BufProps map[string]bProperty     `yaml:"properties"`
	BufConns map[string][]bConnection `yaml:"connections"`
	nesting
}

// Temporal buffer type for .yaml parsing purposes; represents
// nodes-field of template-file
type bNode struct {
	BufLabels []string                 `yaml:"labels"`
	BufProps  map[string]bProperty     `yaml:"properties"`
	BufConns  map[string][]bConnection `yaml:"connections"`
	nesting
}

// Temporal buffer type for .yaml parsing purposes; represents
// edges-field of template-file
type bEdge struct {
	BufProps map[string]bProperty `yaml:"properties"`
	nesting
}

// Temporal buffer type for .yaml parsing purposes; represents
// subfileds of properties-field within "nodes" and "edges"
type bProperty struct {
	BufType   string        `yaml:"type"`
	BufRestrs bRestrictions `yaml:"restrictions"`
	nesting
}

// Temporal buffer type for .yaml parsing purposes; represents
// restriction-subfiled of property-field within "nodes", "labels"
// and "edges"
type bRestrictions struct {
	BufValueRestr     []string `yaml:"values"`
	BufRegexpRestr    []string `yaml:"regexps"`
	BufKeyValueRestr  []string `yaml:"key_values"`
	BufKeyRegexpRestr []string `yaml:"key_regexps"`
	nesting
}

// Temporal buffer type for .yaml parsing purposes; represents
// subfileds of connections-field within "nodes" and "labels"
type bConnection struct {
	BufEdge  string `yaml:"edge"`
	BufRatio struct {
		Min int `yaml:"min"`
		Max int `yaml:"max"`
	} `yaml:"ratio"`
	nesting
}
