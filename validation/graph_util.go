package validation

import (
	"fmt"
	"reflect"
)

// ----------------------------- INSERT/DELETE ------------------------------ //

// Inserts n node to the g graph; omits insertation if identical node is already presented
// within g graph; returns true at success
func (g graphHolder) insertNode(n Node) bool {
	m := g.searchMainNodeMatchWide(n)
	if m != nil {
		return false
	}

	m = &graphHolderNode{
		node:   n,
		childs: make(map[*graphHolderNode][]Edge),
	}
	g[m] = struct{}{}

	return true
}

// Removes n node from the g graph (including all references within another nodes); returns
// true at success
func (g graphHolder) deleteNode(n Node) bool {
	// searches for references
	for m := range g {
		for s := range m.childs {
			if !isEqualNode(s.node, n) {
				continue
			}
			delete(m.childs, s)
		}
	}
	// searches for node itself
	for m := range g {
		if !isEqualNode(m.node, n) {
			continue
		}
		delete(g, m)
		m = nil
		return true
	}

	return false
}

// Inserts the whole tr triplet within g graph; returns true at success
//
// In some cases has a little bit different semanthic:
//   - already presented edge - cancels insertation
//   - already presented main node and subject node - connects them with
//     edge (adding edge)
//   - nil subject node and nil edge - inserts main node as single node
//     (adding node)
//   - nil main node and nil edge - inserts subject node as single node
//     (adding node)
//   - nil edge - inserts main node and subject node as single nodes (adding
//     nodes)
func (g graphHolder) insertTriplet(tr Triplet) bool {
	m, s, e := g.searchForMatches(tr)

	if e != nil {
		// it means that there is identical triplet within g exists
		return false
	}

	res := false

	if m == nil && tr.Main() != nil {
		// it means that there is no any identical node within g exists
		m = &graphHolderNode{
			node:   tr.Main(),
			childs: make(map[*graphHolderNode][]Edge),
		}
		g[m] = struct{}{}
		res = true
	}

	if s == nil && tr.Subj() != nil {
		// it means that there is no any identical node within g exists
		s = &graphHolderNode{
			node:   tr.Subj(),
			childs: make(map[*graphHolderNode][]Edge),
		}
		g[s] = struct{}{}
		res = true
	}

	if m != nil && s != nil && tr.Edge() != nil {
		m.childs[s] = append(m.childs[s], tr.Edge())
		res = true
	}

	return res
}

// Removes edge, described by tr triplet, from the g graph; returns true at success
func (g graphHolder) deleteEdge(tr Triplet) bool {
	m, s, e := g.searchForMatches(tr)

	if e != nil {
		// it means that there is identical triplet within g exists
		for i, edge := range m.childs[s] {
			if isEqualEdge(e, edge) {
				m.childs[s] = append(m.childs[s][:i], m.childs[s][i+1:]...)
				return true
			}
		}
	}

	return false
}

// --------------------------------- SEARCH --------------------------------- //

// Searches for triplets described by filters type-names and returns them
//
// WARNING: filters should have at MOST 3 args, where 1st arg is main node type,
// 2nd - subject node type and 3rd - edge type; any excess args will be OMITTED
func (g graphHolder) findTriplets(filters ...string) []Triplet {
	fs := [3]string{}
	// [0] - mainType filter
	// [1] - subjType filter
	// [2] - edgeType filter
	// indexes above 2 will be omitted
	for i := 0; i < len(filters) && i < 3; i++ {
		fs[i] = filters[i]
	}

	res := make([]Triplet, 0)
	for m := range g {
		main := m.node
		if fs[0] != "" && fs[0] != main.GetNodeType() {
			continue
		}
		for s, es := range m.childs {
			subj := s.node
			if fs[1] != "" && fs[1] != subj.GetNodeType() {
				continue
			}
			for _, edge := range es {
				if fs[2] != "" && fs[2] != edge.GetEdgeType() {
					continue
				}
				res = append(res, NewTriplet(main, subj, edge))
			}
		}
	}
	return res
}

// Searches for triplets where n node is main node and returns them
func (g graphHolder) getGroupedTriplets(n Node) []Triplet {
	res := make([]Triplet, 0)
	for m := range g {
		if !isEqualNode(m.node, n) {
			continue
		}
		for s, es := range m.childs {
			for _, e := range es {
				res = append(res, NewTriplet(m.node, s.node, e))
			}
		}
		break
	}
	return res
}

// Searches for nodes described by filter node-type-name and returns them
func (g graphHolder) findNodes(filter string) []Node {
	res := make([]Node, 0)
	for n := range g {
		main := n.node
		if filter != "" && filter != main.GetNodeType() {
			continue
		}
		res = append(res, main)
	}
	return res
}

// Searches for triplets where n node is main node and returns them as
// duplets (edge + subject node)
func (g graphHolder) getNodeChilds(n Node) []Duplet {
	res := make([]Duplet, 0)
	for m := range g {
		if !isEqualNode(m.node, n) {
			continue
		}
		for s, es := range m.childs {
			for _, e := range es {
				res = append(res, NewDuplet(s.node, e))
			}
		}
		break
	}
	return res
}

// -------------------------------- HELPERS --------------------------------- //

// Searches first (an the only one by desing) occurrence of tr triplet; returns non-nil
// results at succcess (first - main node, second - subject node, third - edge)
func (g graphHolder) searchForMatches(tr Triplet) (*graphHolderNode, *graphHolderNode, Edge) {
	var (
		matchMain *graphHolderNode
		matchSubj *graphHolderNode
		matchEdge Edge
	)

	if tr.Main() == nil {
		return nil, nil, nil
	}
	matchMain = g.searchMainNodeMatchWide(tr.Main())
	if matchMain != nil {
		if tr.Subj() == nil {
			return matchMain, nil, nil
		}
		matchSubj = searchSubjNodeMatchDeep(matchMain, tr.Subj())
		if matchSubj != nil {
			if tr.Edge() == nil {
				return matchMain, matchSubj, nil
			}
			matchEdge = searchEdgeMatchDeep(matchMain, matchSubj, tr.Edge())
			return matchMain, matchSubj, matchEdge // happy path - full match
		}
		matchSubj = g.searchSubjNodeMatchWide(tr.Subj())
	}
	return matchMain, matchSubj, nil // matchSubj may be nil
}

func (g graphHolder) searchMainNodeMatchWide(m Node) *graphHolderNode {
	for n := range g {
		if isEqualNode(n.node, m) {
			return n
		}
	}
	return nil
}

// Searches exact node withing g graph using s subject node; returns non-nil
// graphHolderNode-pointer at succcess
func (g graphHolder) searchSubjNodeMatchWide(s Node) *graphHolderNode {
	for n := range g {
		if isEqualNode(n.node, s) {
			return n
		}
	}
	return nil
}

// Searches exact node withing m graph-node's information using s subject node;
// returns non-nil graphHolderNode-pointer at succcess
func searchSubjNodeMatchDeep(m *graphHolderNode, s Node) *graphHolderNode {
	for n := range m.childs {
		if isEqualNode(n.node, s) {
			return n
		}
	}
	return nil
}

// Searches exact edge withing m graph-node's information using s subject node and
// e Edge; returns non-nil Edge-interface at succcess
func searchEdgeMatchDeep(m, s *graphHolderNode, e Edge) Edge {
	for _, edge := range m.childs[s] {
		if isEqualEdge(edge, e) {
			return e
		}
	}
	return nil
}

// Returns true if n1 and n2 nodes are equal
func isEqualNode(n1, n2 Node) bool {
	if n1.GetNodeType() != n2.GetNodeType() {
		return false
	}
	if len(n1.GetKeys()) != len(n2.GetKeys()) {
		return false
	}
	for _, k := range n1.GetKeys() {
		res1, ok1 := n1.GetProp(k)
		res2, ok2 := n2.GetProp(k)
		if !ok1 && !ok2 {
			return false
		}

		if ok, err := tryCompareAsComplexDataType(res1, res2); ok {
			continue
		} else if err == nil {
			return false
		}
		if res1 != res2 {
			return false
		}
	}
	return true
}

// Returns true if e1 and e2 edges are equal
func isEqualEdge(e1, e2 Edge) bool {
	if e1.GetEdgeType() != e2.GetEdgeType() {
		return false
	}
	if len(e1.GetKeys()) != len(e2.GetKeys()) {
		return false
	}
	for _, k := range e1.GetKeys() {
		res1, ok1 := e1.GetProp(k)
		res2, ok2 := e2.GetProp(k)
		if !ok1 && !ok2 {
			return false
		}

		if ok, err := tryCompareAsComplexDataType(res1, res2); ok {
			continue
		} else if err == nil {
			return false
		}
		if res1 != res2 {
			return false
		}

	}
	return true
}

// Tries to assert equality of values as their underlying values are maps or
// slices/arrays; panics if both values has different types (map and slice for
// example); returns non-nil error if p1 and p2 underlying values are not maps or
// slices
func tryCompareAsComplexDataType(p1, p2 interface{}) (bool, error) {
	v1 := reflect.ValueOf(p1)
	v2 := reflect.ValueOf(p2)
	if v1.Kind() != v2.Kind() {
		return false, fmt.Errorf("values has not equal data types")
	}
	switch v1.Kind() {
	case reflect.Map:
		return isEqualMapProperty(v1, v2), nil
	case reflect.Slice, reflect.Array:
		return isEqualArrayProperty(v1, v2), nil
	default:
		return false, fmt.Errorf("values are not maps or arrays/slices")
	}
}

// Returns true if given p1 and p2 properties (wrapped in reflect.Value-structs) are
// equal; panics if p1 and p2 underlying values are NOT maps
func isEqualMapProperty(p1, p2 reflect.Value) bool {
	if p1.Len() != p2.Len() {
		return false
	}
	for iter1 := p1.MapRange(); iter1.Next(); {

		key1 := iter1.Key()
		checkVal2 := p2.MapIndex(key1)
		if !checkVal2.IsValid() {
			return false
		}
		val1 := iter1.Value().Interface()
		val2 := checkVal2.Interface()
		if val1 != val2 {
			return false
		}
	}
	return true
}

// Returns true if given p1 and p2 properties (wrapped in reflect.Value-structs) are
// equal; panics if p1 and p2 underlying values are NOT arrays or slices
func isEqualArrayProperty(p1, p2 reflect.Value) bool {
	if p1.Len() != p2.Len() {
		return false
	}
	for i := 0; i < p1.Len(); i++ {
		val1 := p1.Index(i).Interface()
		val2 := p2.Index(i).Interface()
		if val1 != val2 {
			return false
		}
	}
	return true
}

// Pretty-prints graph info; used only for debug
func (g graphHolder) debugPrint() {
	mId := 0
	mp := make(map[*graphHolderNode]int)
	for m := range g {
		mId += 1
		mp[m] = mId
	}
	lines := make([]string, 0)

	left, right := 0, 0
	for m := range g {
		left += 1
		right = left
		sLastInd := 0

		sCount := 0
		{
			node := (m.node).(nodeHolder)
			l := fmt.Sprintf("(%d) (%q) <<%+v>>\n", mp[m], node.typ, node.props)
			lines = append(lines, l)
		}
		for s, e := range m.childs {
			sLastInd = right
			innerLeft := sLastInd + 1
			innerRight := sLastInd
			right += 1

			sCount += 1
			{
				l := fmt.Sprintf(" '--(%d)\n", mp[s])
				lines = append(lines, l)
			}
			for _, e := range e {
				right += 1
				innerRight += 1
				{
					edge := (e).(edgeHolder)
					l := fmt.Sprintf("     '--[%q] <<%+v>>\n", edge.typ, edge.props)
					lines = append(lines, l)
				}
			}
			if innerRight-innerLeft >= 1 {
				for i := innerLeft; i < innerRight; i++ {
					lines[i] = lines[i][:5] + "|" + lines[i][6:]
				}
			}
		}
		if sLastInd-left > 1 && sCount > 1 {
			for i := left; i < sLastInd; i++ {
				lines[i] = " |" + lines[i][2:]
			}
		}
		lines = append(lines,
			"----------------------------------------------\n")
		left += right - left + 1
	}

	fmt.Print("----------------------------------------------\n")
	for _, line := range lines {
		fmt.Print(line)
	}
}
