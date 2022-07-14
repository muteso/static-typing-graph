# **S**tatic **T**yping for **G**raph data models (**STG**)

## Definition

This package provides instruments for convinient wrapping of any possible data type to "staticly typed" nodes and edges of graph.

This package should **not** be treated as ORM because it does not provide any kind of "adapters" to existing graph data bases (like neo4j DB or dgraph) - the main purpose of this package is to wrap standart go types into graph entities and restrict them like they have static types.

## How to use

The first (and the most difficult) step you should execute - is to define graph data model which data types will be then used for wrapping your in-programm data.

After parsing of this definition, this graph data model will be covered by Validator-interface, which then will be used for restriction of wrapped data (i.e. ensuring static typing).
Definition of graph data model consists from 3 types - nodes, edges and labels (+ 2 implicit types - property and connection). This definitin **must** to be written in **yaml**-notation - you may do this by writing yaml-file or in-programm string.

So, we have 3 types:
- nodes - the main data holders; have type name definition, label-embedding definition, properties definition (remember of implicit types) and definition of connections with other node-types (remember of implicit types),
- edges - secondary data holders which used **only** within node-connection definitions; have type name definition and properties definition,
- labels - type which used **only** for embedding data within node types; thus labels provide easy way of reusing some chunks of definitions between any amount of node-types; have the same definitions as nodes excluding label-embedding (to avoid bad graph-design decisions, endless data-nesting and property-name conflicts); **IMPORTANT**: if node-type and label-type have the "same" definitions - node-type's definion **overwrites** label-type's definion,

... and 2 implicit types:
- connections
- properties
  - int
  - float
  - bool
  - string
  - datetime
  - array
  - map

This graph defenition reference looks like this:
```
labels: # may be empty
  <type name>:
    properties: # may be omitted
      <property name>:
        type: <data type>
        restrictions:
          values: # may be omitted
            - <first variant>
            - <etc...>
          regexps: # may be omitted
            - <first variant>
            - <etc...>
          key_values: # can be used only if type of property is 'map'; may be omitted
            - <first variant>
            - <etc...>
          key_regexps: # can be used only if type of property is 'map'; may be omitted
            - <first variant>
            - <etc...>
    connections: # may be omitted
      <label type name which connects with this label type>:
        - edge: <edge type name which is used to connect this label with label mentoined above>
          ratio: 
            min: <min amount of unique instances of nodes, which contains label mentoined above, connected with a single instance of node, which contains this label>
            max: <max amount of unique instances of nodes, which contains label mentoined above, connected with a single instance of node, which contains this label>
nodes: # may be empty
  <type name>:
    labels: # may be omitted
      - <label name, which is defined above>
      - <etc...>
    properties: # may be omitted
      <property name>:
        type: <data type>
        restrictions:
          values: # may be omitted
            - <first variant>
            - <etc...>
          regexps: # may be omitted
            - <first variant>
            - <etc...>
          key_values: # can be used only if type of property is 'map'; may be omitted
            - <first variant>
            - <etc...>
          key_regexps: # can be used only if type of property is 'map'; may be omitted
            - <first variant>
            - <etc...>
    connections: # may be omitted
      <node type name which connects with this node type>:
        - edge: <edge type name which is used to connect this node with node mentoined above>
          ratio: 
            min: <min amount of unique instances of node mentoined above connected with a single instance of this node>
            max: <max amount of unique instances of node mentoined above connected with a single instance of this node>
edges: # may be empty
  <type name>:
    properties: # may be omitted
      <property name>:
        type: <data type>
        restrictions:
          values: # may be omitted
            - <first variant>
            - <etc...>
          regexps: # may be omitted
            - <first variant>
            - <etc...>
          key_values: # can be used only if type of property is 'map'; may be omitted
            - <first variant>
            - <etc...>
          key_regexps: # can be used only if type of property is 'map'; may be omitted
            - <first variant>
            - <etc...>
```