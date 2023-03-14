package stg

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var templ Validator
var testTime time.Time

func init() {
	file, err := os.Open("template_example.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	templ, err = ParseTemplate(file)
	if err != nil {
		panic(err)
	}

	testTime, _ = time.Parse(time.RFC3339, "1111-11-11T11:11:11Z")
}

func TestValidateNode(t *testing.T) {
	person := NewNode("Person", map[string]interface{}{
		"name":    "Jora",
		"birth":   testTime,
		"merried": true,
		"age":     22.7,
		"money":   34,
		"things":  []string{"other thing"},
		"adresses": map[string]string{
			"street 1": "house  ",
		},
	})
	if ok, err := Validate(templ, person); !ok {
		t.Error("Is NOT valid: Node-interface -> " + err.Error())
	} else {
		fmt.Println("Is valid: Node-interface")
	}
}

func TestValidateEdge(t *testing.T) {
	friend := NewEdge("friend", map[string]interface{}{
		"since": testTime,
	})
	if ok, err := Validate(templ, friend); !ok {
		t.Error("Is NOT valid: Edge-interface -> " + err.Error())
	} else {
		fmt.Println("Is valid: Edge-interface")
	}
}

func TestValidateDuplet(t *testing.T) {
	person := NewNode("Person", map[string]interface{}{
		"name":    "Jora",
		"birth":   testTime,
		"merried": true,
		"age":     22.7,
		"money":   34,
		"things":  []string{"thing", "other thing"},
		"adresses": map[string]string{
			"street 1": "house  ",
		},
	})
	friend := NewEdge("friend", map[string]interface{}{
		"since": testTime,
	})
	duplet := NewDuplet(person, friend)
	if ok, err := Validate(templ, duplet); !ok {
		t.Error("Is NOT valid: Duplet-interface -> " + err.Error())
	} else {
		fmt.Println("Is valid: Duplet-interface")
	}
}

func TestValidateTriplet(t *testing.T) {
	person := NewNode("Person", map[string]interface{}{
		"name":    "Jora",
		"birth":   testTime,
		"merried": true,
		"age":     22.7,
		"money":   34,
		"things":  []string{"thing"},
		"adresses": map[string]string{
			"street 1": "house  ",
		},
	})
	friend := NewEdge("friend", map[string]interface{}{
		"since": testTime,
	})
	triplet := NewTriplet(person, person, friend)
	if ok, err := Validate(templ, triplet); !ok {
		t.Error("Is NOT valid: Triplet-interface -> " + err.Error())
	} else {
		fmt.Println("Is valid: Triplet-interface")
	}
}

func TestValidateUnknownStruct(t *testing.T) {
	{
		type Person struct {
			Name     string
			Birth    time.Time
			Merried  bool
			Age      float64
			Money    int
			Things   []string
			Adresses map[string]string
		}
		person := Person{
			"Jora",
			testTime,
			true,
			22.7,
			34,
			[]string{"thing"},
			map[string]string{
				"street 1": "house 12",
			},
		}
		if ok, err := Validate(templ, person); !ok {
			t.Error("Is NOT valid: Person-struct -> " + err.Error())
		} else {
			fmt.Println("Is valid: Person-struct")
		}
	}
	{
		type PersonStruct struct {
			P        string
			Name     string
			Birth    time.Time
			Merried  bool
			Age      float64
			Money    int
			Things   []string
			Adresses map[string]string
		}
		person := PersonStruct{
			"Person",
			"Jora",
			testTime,
			true,
			22.7,
			34,
			[]string{"thing"},
			map[string]string{
				"street 1": "house  ",
			},
		}
		if ok, err := Validate(templ, person); !ok {
			t.Error("Is NOT valid: PersonStruct-struct with Person-value -> " + err.Error())
		} else {
			fmt.Println("Is valid: PersonStruct-struct with Person-value")
		}
	}
}

func TestValidateUnknownMap(t *testing.T) {
	{
		type Person map[string]interface{}
		person := Person{
			"name":    "Jora",
			"birth":   testTime,
			"merried": true,
			"age":     22.7,
			"money":   34,
			"things":  []string{"thing"},
			"adresses": map[string]string{
				"street 1": "house  ",
			},
		}
		if ok, err := Validate(templ, person); !ok {
			t.Error("Is NOT valid: Person-map based type -> " + err.Error())
		} else {
			fmt.Println("Is valid: Person-map based type")
		}
	}
	{
		type PersonMap map[string]interface{}
		person := PersonMap{
			"p":       "Person",
			"name":    "Jora",
			"birth":   testTime,
			"merried": true,
			"age":     22.7,
			"money":   34,
			"things":  []string{"thing", "other thing"},
			"adresses": map[string]string{
				"street 1": "house  ",
			},
		}
		if ok, err := Validate(templ, person); !ok {
			t.Error("Is NOT valid: PersonMap-map based type with Person-value -> " + err.Error())
		} else {
			fmt.Println("Is valid: PersonMap-map based type with Person-value")
		}
	}
}

func TestValidateGraph(t *testing.T) {
	// refactor
	n1 := NewNode("Person", map[string]interface{}{
		"name":    "Jora",
		"birth":   testTime,
		"merried": true,
		"age":     22.7,
		"money":   34,
		"things":  []string{"thing"},
		"adresses": map[string]string{
			"street 1": "house 1",
		}})
	n2 := NewNode("Person", map[string]interface{}{
		"name":    "Jora",
		"birth":   testTime,
		"merried": true,
		"age":     22.7,
		"money":   34,
		"things":  []string{"thing", "other thing"},
		"adresses": map[string]string{
			"street 1": "house 2",
		}})
	n3 := NewNode("Person", map[string]interface{}{
		"name":    "Jora",
		"birth":   testTime,
		"merried": true,
		"age":     22.7,
		"money":   34,
		"things":  []string{"thing"},
		"adresses": map[string]string{
			"street 1": "house 3",
		}})
	e1 := NewEdge("friend", map[string]interface{}{
		"since": testTime.Add(100),
	})
	e2 := NewEdge("friend", map[string]interface{}{
		"since": testTime.Add(200),
	})
	tr1 := NewTriplet(n1, n2, e1) // n1 -e1-> n2
	tr2 := NewTriplet(n1, n2, e2) // n1 -e2-> n2
	tr3 := NewTriplet(n2, n3, e1) // n2 -e1-> n3
	tr4 := NewTriplet(n2, n1, e2) // n2 -e2-> n1
	tr5 := NewTriplet(n3, n1, e1) // n3 -e1-> n1

	gr := NewGraph(nil, []Triplet{tr1, tr2, tr3, tr4, tr5}...)
	if ok, err := Validate(templ, gr); !ok {
		t.Error("Is NOT valid: graph -> " + err.Error())
	} else {
		fmt.Println("Is valid: graph")
	}

	// res := gr.GetTripletsByType("Person")
	// fmt.Printf("%v\n", res)

	// gr.RemoveEdge(tr1)
	// res1 := gr.GetTripletsByType("Person")
	// fmt.Printf("%v\n", res1)

	// res2 := gr.GetNodes()
	// fmt.Printf("%v\n", res2)

	// res3 := gr.GetNodeChilds(n1)
	// fmt.Printf("%v\n", res3)
}
