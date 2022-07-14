package parser

import (
	"stg/template"
	"sync"
	"testing"
)

func TestNewContext(t *testing.T) {
	temp := newContext()
	if temp.res == nil ||
		temp.ls == nil ||
		temp.nls == nil ||
		temp.lcn == nil ||
		temp.errs == nil ||
		temp.Mutex == nil {
		t.Error("Test-case is failed")
	}
}

func TestAppendErr(t *testing.T) {
	temp := &context{
		errs: []parseError{},
	}
	temp.appendErr(parseError{
		Loc: "test0",
		Msg: "test1",
	})
	exp := temp.errs[0]
	if exp.Loc != "test0" || exp.Msg != "test1" {
		t.Error("Test-case is failed")
	}
}

func TestBuildErr(t *testing.T) {
	temp := &context{
		errs: []parseError{
			{
				Loc: "test0",
				Msg: "test1",
			},
			{
				Loc: "test2",
				Msg: "test3",
			},
		},
	}
	exp := temp.buildErr()
	if exp.Error() != "test0 >> test1\ntest2 >> test3\n" {
		t.Error("Successive test case is failed")
	}
	temp = &context{
		errs: []parseError{},
	}
	exp = temp.buildErr()
	if exp != nil {
		t.Error("Unsuccessive test case is failed")
	}
}

func TestNode(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{
				"Test": {},
			},
		},
		Mutex: new(sync.Mutex),
	}
	exp := temp.res.Nodes["Test"]
	if temp.node("Test") != exp {
		t.Error("Successive test-case is failed")
	}
	if temp.node("TestFail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.node("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res:   &template.TemplateHolder{},
		Mutex: new(sync.Mutex),
	}
	if temp.node("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.node("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestNodeProp(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{
				"Test": {
					Props: map[string]*template.TProperty{
						"test": {},
					},
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	exp := temp.res.Nodes["Test"].Props["test"]
	if temp.nodeProp("Test", "test") != exp {
		t.Error("Successive test-case is failed")
	}
	if temp.nodeProp("Test", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	if temp.nodeProp("TestFail", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{
				"Test": {
					Props: map[string]*template.TProperty{},
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.nodeProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{
				"Test": {},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.nodeProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.nodeProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.nodeProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestNodeLabels(t *testing.T) {
	temp := &context{
		nls: map[string][]string{
			"Test": {
				"test", // necessary to comparsion
			},
		},
		Mutex: new(sync.Mutex),
	}
	exp := temp.nls["Test"][0]
	if temp.nodeLabels("Test")[0] != exp {
		t.Error("Successive test-case is failed")
	}
	if temp.nodeLabels("TestFail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		nls:   map[string][]string{},
		Mutex: new(sync.Mutex),
	}
	if temp.nodeLabels("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.nodeLabels("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestLabel(t *testing.T) {
	temp := &context{
		ls: map[string]*cLabel{
			"Test": {},
		},
		Mutex: new(sync.Mutex),
	}
	exp := temp.ls["Test"]
	if temp.label("Test") != exp {
		t.Error("Successive test-case is failed")
	}
	if temp.label("TestFail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		ls:    map[string]*cLabel{},
		Mutex: new(sync.Mutex),
	}
	if temp.label("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.label("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestLabelProp(t *testing.T) {
	temp := &context{
		ls: map[string]*cLabel{
			"Test": {
				props: map[string]*template.TProperty{
					"test": {},
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	exp := temp.ls["Test"].props["test"]
	if temp.labelProp("Test", "test") != exp {
		t.Error("Successive test-case is failed")
	}
	if temp.labelProp("Test", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	if temp.labelProp("TestFail", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		ls: map[string]*cLabel{
			"Test": {
				props: map[string]*template.TProperty{},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.labelProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		ls: map[string]*cLabel{
			"Test": {},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.labelProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		ls:    map[string]*cLabel{},
		Mutex: new(sync.Mutex),
	}
	if temp.labelProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.labelProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestLabelNodes(t *testing.T) {
	temp := &context{
		ls: map[string]*cLabel{
			"Test": {
				nodes: map[string]*template.TNode{
					"test": {}, // necessary to comparsion
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	exp := temp.ls["Test"].nodes["test"]
	if temp.labelNodes("Test")["test"] != exp {
		t.Error("Successive test-case is failed")
	}
	if temp.labelNodes("TestFail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		ls: map[string]*cLabel{
			"Test": {},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.labelNodes("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		ls:    map[string]*cLabel{},
		Mutex: new(sync.Mutex),
	}
	if temp.labelNodes("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.labelNodes("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestEdge(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{
				"Test": {},
			},
		},
		Mutex: new(sync.Mutex),
	}
	exp := temp.res.Edges["Test"]
	if temp.edge("Test") != exp {
		t.Error("Successive test-case is failed")
	}
	if temp.edge("TestFail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.edge("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res:   &template.TemplateHolder{},
		Mutex: new(sync.Mutex),
	}
	if temp.edge("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.edge("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestEdgeProp(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{
				"Test": {
					Props: map[string]*template.TProperty{
						"test": {},
					},
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	exp := temp.res.Edges["Test"].Props["test"]
	if temp.edgeProp("Test", "test") != exp {
		t.Error("Successive test-case is failed")
	}
	if temp.edgeProp("Test", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	if temp.edgeProp("TestFail", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{
				"Test": {
					Props: map[string]*template.TProperty{},
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.edgeProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{
				"Test": {},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.edgeProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.edgeProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		res:   &template.TemplateHolder{},
		Mutex: new(sync.Mutex),
	}
	if temp.edgeProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.edgeProp("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestLabelConnsByMain(t *testing.T) {
	temp := &context{
		lcn: map[string]map[string]map[string]*cLConnection{
			"Test": {},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.labelConnsByMain("Test") == nil {
		t.Error("Successive test-case is failed")
	}
	if temp.labelConnsByMain("TestFail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	// there is no point to test all map nesting variants because even
	// if you try to fetch any value from nil map you'll get default
	// value of this map - in out case its nil
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.labelConnsByMain("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestLabelConnsByMainSubj(t *testing.T) {
	temp := &context{
		lcn: map[string]map[string]map[string]*cLConnection{
			"Test": {
				"test": {},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.labelConnsByMainSubj("Test", "test") == nil {
		t.Error("Successive test-case is failed")
	}
	if temp.labelConnsByMainSubj("Test", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	if temp.labelConnsByMainSubj("TestFail", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	// there is no point to test all map nesting variants because even
	// if you try to fetch any value from nil map you'll get default
	// value of this map - in out case its nil
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.labelConnsByMainSubj("Test", "test") != nil {
		t.Error("Successive test-case is failed")
	}
}

func TestLabelConn(t *testing.T) {
	temp := &context{
		lcn: map[string]map[string]map[string]*cLConnection{
			"Test": {
				"test": {
					"t": {},
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.labelConn("Test", "test", "t") == nil {
		t.Error("Successive test-case is failed")
	}
	if temp.labelConn("Test", "testfail", "tf") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	if temp.labelConn("Test", "test", "tf") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	if temp.labelConn("TestFail", "testfail", "tf") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	// there is no point to test all map nesting variants because even
	// if you try to fetch any value from nil map you'll get default
	// value of this map - in out case its nil
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.labelConn("Test", "test", "t") != nil {
		t.Error("Successive test-case is failed")
	}
}

func TestNodeConnsByMain(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Conns: map[string]map[string]map[string]*template.TConnection{
				"Test": {},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.nodeConnsByMain("Test") == nil {
		t.Error("Successive test-case is failed")
	}
	if temp.nodeConnsByMain("TestFail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	// there is no point to test all map nesting variants because even
	// if you try to fetch any value from nil map you'll get default
	// value of this map - in out case its nil
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.nodeConnsByMain("Test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestNodeConnsByMainSubj(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Conns: map[string]map[string]map[string]*template.TConnection{
				"Test": {
					"test": {},
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.nodeConnsByMainSubj("Test", "test") == nil {
		t.Error("Successive test-case is failed")
	}
	if temp.nodeConnsByMainSubj("Test", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	if temp.nodeConnsByMainSubj("TestFail", "testfail") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	// there is no point to test all map nesting variants because even
	// if you try to fetch any value from nil map you'll get default
	// value of this map - in out case its nil
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.nodeConnsByMainSubj("Test", "test") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestNodeConn(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Conns: map[string]map[string]map[string]*template.TConnection{
				"Test": {
					"test": {
						"t": {},
					},
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	if temp.nodeConn("Test", "test", "t") == nil {
		t.Error("Successive test-case is failed")
	}
	if temp.nodeConn("Test", "testfail", "tf") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	if temp.nodeConn("Test", "test", "tf") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	if temp.nodeConn("TestFail", "testfail", "tf") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
	// there is no point to test all map nesting variants because even
	// if you try to fetch any value from nil map you'll get default
	// value of this map - in out case its nil
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	if temp.nodeConn("Test", "test", "t") != nil {
		t.Error("Unsuccessive test-case is failed")
	}
}

func TestSetNode(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{},
		},
		Mutex: new(sync.Mutex),
	}
	obj := &template.TNode{}
	temp.setNode("Test", obj)
	if temp.res.Nodes["Test"] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		res:   &template.TemplateHolder{},
		Mutex: new(sync.Mutex),
	}
	temp.setNode("Test", obj)
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setNode("Test", obj)
}

func TestSetNodeProp(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{
				"Test": {
					Props: map[string]*template.TProperty{},
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	obj := &template.TProperty{}
	temp.setNodeProp("Test", "test", obj)
	if temp.res.Nodes["Test"].Props["test"] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{
				"Test": {},
			},
		},
		Mutex: new(sync.Mutex),
	}
	temp.setNodeProp("Test", "test", obj)
	temp = &context{
		res: &template.TemplateHolder{
			Nodes: map[string]*template.TNode{},
		},
		Mutex: new(sync.Mutex),
	}
	temp.setNodeProp("Test", "test", obj)
	temp = &context{
		res:   &template.TemplateHolder{},
		Mutex: new(sync.Mutex),
	}
	temp.setNodeProp("Test", "test", obj)
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setNodeProp("Test", "test", obj)
}

func TestSetNodeLabel(t *testing.T) {
	temp := &context{
		nls: map[string][]string{
			"Test": {},
		},
		Mutex: new(sync.Mutex),
	}
	obj := "test"
	temp.setNodeLabel("Test", obj)
	if temp.nls["Test"][0] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		nls:   map[string][]string{},
		Mutex: new(sync.Mutex),
	}
	temp.setNodeLabel("Test", obj)
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setNodeLabel("Test", obj)
}

func TestSetLabel(t *testing.T) {
	temp := &context{
		ls:    map[string]*cLabel{},
		Mutex: new(sync.Mutex),
	}
	obj := &cLabel{}
	temp.setLabel("Test", obj)
	if temp.ls["Test"] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setLabel("Test", obj)
}

func TestSetLabelProp(t *testing.T) {
	temp := &context{
		ls: map[string]*cLabel{
			"Test": {
				props: make(map[string]*template.TProperty),
			},
		},
		Mutex: new(sync.Mutex),
	}
	obj := &template.TProperty{}
	temp.setLabelProp("Test", "test", obj)
	if temp.ls["Test"].props["test"] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		ls: map[string]*cLabel{
			"Test": {},
		},
		Mutex: new(sync.Mutex),
	}
	temp.setLabelProp("Test", "test", obj)
	temp = &context{
		ls:    map[string]*cLabel{},
		Mutex: new(sync.Mutex),
	}
	temp.setLabelProp("Test", "test", obj)
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setLabelProp("Test", "test", obj)
}

func TestSetLabelNode(t *testing.T) {
	temp := &context{
		ls: map[string]*cLabel{
			"Test": {
				nodes: make(map[string]*template.TNode),
			},
		},
		Mutex: new(sync.Mutex),
	}
	obj := &template.TNode{}
	temp.setLabelNode("Test", "test", obj)
	if temp.ls["Test"].nodes["test"] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		ls: map[string]*cLabel{
			"Test": {},
		},
		Mutex: new(sync.Mutex),
	}
	temp.setLabelNode("Test", "test", obj)
	temp = &context{
		ls:    map[string]*cLabel{},
		Mutex: new(sync.Mutex),
	}
	temp.setLabelNode("Test", "test", obj)
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setLabelNode("Test", "test", obj)
}

func TestSetEdge(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{},
		},
		Mutex: new(sync.Mutex),
	}
	obj := &template.TEdge{}
	temp.setEdge("Test", obj)
	if temp.res.Edges["Test"] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		res:   &template.TemplateHolder{},
		Mutex: new(sync.Mutex),
	}
	temp.setEdge("Test", obj)
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setEdge("Test", obj)
}

func TestSetEdgeProp(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{
				"Test": {
					Props: make(map[string]*template.TProperty),
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	obj := &template.TProperty{}
	temp.setEdgeProp("Test", "test", obj)
	if temp.res.Edges["Test"].Props["test"] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{
				"Test": {},
			},
		},
		Mutex: new(sync.Mutex),
	}
	temp.setEdgeProp("Test", "test", obj)
	temp = &context{
		res: &template.TemplateHolder{
			Edges: map[string]*template.TEdge{},
		},
		Mutex: new(sync.Mutex),
	}
	temp.setEdgeProp("Test", "test", obj)
	temp = &context{
		res:   &template.TemplateHolder{},
		Mutex: new(sync.Mutex),
	}
	temp.setEdgeProp("Test", "test", obj)
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setEdgeProp("Test", "test", obj)
}

// func TestSetNodeConnsByMain(t *testing.T) {}

// func TestSetNodeConnsByMainSubj(t *testing.T) {}

func TestSetNodeConn(t *testing.T) {
	temp := &context{
		res: &template.TemplateHolder{
			Conns: map[string]map[string]map[string]*template.TConnection{
				"Test": {
					"test": make(map[string]*template.TConnection),
				},
			},
		},
		Mutex: new(sync.Mutex),
	}
	obj := &template.TConnection{}
	temp.setNodeConn("Test", "test", "t", obj)
	if temp.res.Conns["Test"]["test"]["t"] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		res: &template.TemplateHolder{
			Conns: map[string]map[string]map[string]*template.TConnection{
				"Test": {},
			},
		},
		Mutex: new(sync.Mutex),
	}
	temp.setNodeConn("Test", "test", "t", obj)
	temp = &context{
		res: &template.TemplateHolder{
			Conns: map[string]map[string]map[string]*template.TConnection{},
		},
		Mutex: new(sync.Mutex),
	}
	temp.setNodeConn("Test", "test", "t", obj)
	temp = &context{
		res:   &template.TemplateHolder{},
		Mutex: new(sync.Mutex),
	}
	temp.setNodeConn("Test", "test", "t", obj)
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setNodeConn("Test", "test", "t", obj)
}

// func TestSetLabelConnsByMain(t *testing.T) {}

// func TestSetLabelConnsByMainSubj(t *testing.T) {}

func TestSetLabelConn(t *testing.T) {
	temp := &context{
		lcn: map[string]map[string]map[string]*cLConnection{
			"Test": {
				"test": make(map[string]*cLConnection),
			},
		},
		Mutex: new(sync.Mutex),
	}
	obj := &cLConnection{}
	temp.setLabelConn("Test", "test", "t", obj)
	if temp.lcn["Test"]["test"]["t"] != obj {
		t.Error("Successive test-case is failed")
	}
	defer func() {
		if e := recover(); e != nil {
			t.Error("Unsuccessive test-case is failed")
		}
	}()
	temp = &context{
		lcn: map[string]map[string]map[string]*cLConnection{
			"Test": {},
		},
		Mutex: new(sync.Mutex),
	}
	temp.setLabelConn("Test", "test", "t", obj)
	temp = &context{
		lcn:   map[string]map[string]map[string]*cLConnection{},
		Mutex: new(sync.Mutex),
	}
	temp.setLabelConn("Test", "test", "t", obj)
	temp = &context{
		Mutex: new(sync.Mutex),
	}
	temp.setLabelConn("Test", "test", "t", obj)
}
