# **S**tatic **T**yping for **G**raph data models (**STG**)

## Definition

This package provides instruments for convinient wrapping of go data to "staticly typed" graph-entities.

This package should **not** be treated as ORM because it does not provide any kind of "adapters" to existing graph data bases (like neo4j DB or dgraph) - the main purpose of this package is to wrap standart go types into graph entities and restrict them like they have static types (+ some additioanl restriction options).

## How to use

The first (and the most difficult) step you should execute - is to define graph's data types which are all together is called **template**. Template's data types will be then used to validate your in-programm data.

Definition of template includes sets of 3 types - nodes, edges and labels (+ 2 "implicit" types - properties and connections). This definitin **must** be written in **yaml**-notation - you may do this by writing yaml-file or in-programm string.

<details>
  <summary>Explanation of template definition</summary>

So, we have 3 types:
- nodes - the main data holders; have type name definition, label-embedding definition, properties definition (remember of implicit types) and definition of connections with other node-types (remember of implicit types),
- edges - secondary data holders which used **only** within node-connection definitions; have type name definition and properties definition,
- labels - type which used **only** for embedding data within node types; thus labels provide easy way of reusing some chunks of definitions between any amount of node-types; have the same definitions as nodes excluding label-embedding (to avoid endless data-nesting and definition conflicts); **IMPORTANT**: if node-type and label-type have the "same" definitions - node-type's definion **overwrites** label-type's definion,

... and 2 implicit types:
- connections - describes how nodes is interconnected with each other by edges; connection describes only the **outgoing** interconnections so semanticly they should have the main node (which is contain connection definition), the subject node (which is mentoined in the head of connection definition) and the edge (which is mentoined in body of connection definition); have subject node- and edge-type names references and ratio definition; ratio definition has a few rules:
  - min-field describes minimum amount of **unique** subject nodes which may be connected with a **single** main node using edge; may take any value >= 0,
  - max-field describes maximum amount of **unique** subject nodes which may be connected with a **single** main node using edge; may take any value > 0 (because if max-field takes 0 value the whole connection definition just dont make sense) or may take -1, which means positive infinity (or just any amount > 0),
- properties - describes which information nodes and edges can hold; have property name definition, data type definition (will be described little further) and restrictions definition; restrictions definitions contains exact values and regexeps, where property **must** satisfy at least **one** of them; the available data types to choose in data type definition:
  - "primitive" types:
    - int - int equivalent,
    - float - float64 equivalent,
    - bool - bool equivalent,
    - string - string equivalent,
    - datetime - time.Time equivalent,
  - "complex" types:
    - array - slice/array equivalent; values of array **must** be "primitive" types; definition of array type should look like ```array-<value "primitive" type>```; restrictions for arrays are defined only for the **inner** array's values, not arrays itself,
    - map - map equivalent; keys and values of map **must** be "primitive" types; definition of map type should look like ```map-<key "primitive" type>-<value "primitive" type>```; restrictions for maps are defined only for the **inner** map's values (and keys), not maps itself,

This whole graph defenition reference looks like this:
```
labels: # may be omitted
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
nodes:
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
edges:
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

</details>

<details>
  <summary>Example of correct template definition</summary>

```
labels:
  Creature:
    properties:
      name:
        type: string
        restrictions:
          values:
            - Jora
          regexps:
            - ^[A-Za-z][a-z]+$
nodes:
  Person:
    labels: 
      - Creature
    properties:
      birth:
        type: datetime
        restrictions:
          values:
            - 1111-11-11T11:11:11Z
          regexps:
            - 1111-11
      merried:
        type: bool
        restrictions:
          values:
            - true
      age:
        type: float
        restrictions:
          values:
            - 22.7
          regexps:
            - ^\d+\.\d$
      money:
        type: int
        restrictions:
          values:
            - 34
          regexps:
            - ^\d\d$
      things:
        type: array-string
        restrictions:
          values:
            - thing
          regexps:
            - ^thi..
      adresses:
        type: map-string-string
        restrictions:
          regexps:
            - house .*
          key_regexps:
            - street .+
    connections:
      Person:
        - edge: friend
          ratio: 
            min: 0
            max: -1
edges:
  friend:
    properties:
      since:
        type: datetime
```

</details>


After we define template we should "open" it: 
```
template, err := os.Open("template.yaml")
// or
template := strings.NewReader("\nlables:\n\tCreature:\n\t\t...")
```
... and pass it into the ```stg.ParseTemplate``` function:
```
templ, err := stg.ParseTemplate(template)
```
After successful parsing we will get ```stg.Validator```-interface, which can be used to validate your in-programm data (i.e. ensuring static typing). But firstly, this data **must** be wrapped into another convinient **interfaces** by using this functions:
```
// creating stg.Node-interface
node := stg.NewNode("type name", map[string]interface{}{
  "id": 1,
  "prop1": "value 1",
})

// creating stg.Edge-interface
edge := stg.NewEdge("type name", map[string]interface{}{
  "id": 1,
  "prop1": "value 1",
})

// creating stg.Triplet-interface
triplet := NewTriplet(node, node, edge)

// creating stg.Duplet-interface
duplet := stg.NewDuplet(node, edge)

// creating stg.Graph-interface
graph := stg.NewGraph([]Node{
  node,
}, []stg.Triplet{
  triplet,
}...)
```
And, finally, you can use previously obtained ```stg.Validator``` to validate any value "created" above by using only **one** function - ```stg.Validate```:
```
okNode, nodeError := stg.Validate(templ, node)
...
okGraph, graphError := stg.Validate(templ, graph)
```
As simple as it looks!

## Limitations
Some of the limitations were described above, but it's more convinient to enumerate and repeat all of them here:
  - maps can't nest within each other (which should be handled by making a new node/edge that contains nested map etc.)
  - arrays can't nest within each other (for the same reason as above)
  - maps can't nest within arrays and vise versa (for the same reason as above)
  - map's and array's values (and keys in case of maps) can contain only the same "primitive" (int, string, etc) data types (for reasons of go data types compatibility)
  - labels can't nest within each other (which may lead to endless data-nesting and implicit property- and connection-definitions which may cause hard-to-find definition conflicts)
  - Graph-interface can't contain 2 or more identical nodes (only one instance of node will be created)
  - Graph-interface can't contain 2 or more identical edges with identical directions between one unique pair of nodes (only one instance of edge will be created)