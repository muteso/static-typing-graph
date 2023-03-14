package parser

import (
	"fmt"
	"strings"
	"testing"
)

const (
	file = `
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
                        - street .+
                    key_regexps:
                        - house.*
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
`
	errFile = `
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
                        - 1111-11-11T11:11:11 # error - missing "Z" on end of value
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
                        - street .+
                    key_regexps:
                        - house.*
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
`
	strErr = "template | nodes | Person | properties | birth | restrictions | values | 1 >> restriction \"1111-11-11T11:11:11\" doesn't match \"datetime\" data type\n"
)

func TestParseTemplate(t *testing.T) {
	temp := strings.NewReader(file)
	res, _ := ParseTemplate(temp)
	if res == nil { // refactor - to proper analysis of temp
		t.Error("Successive test case is failed")
	}
	fmt.Println(res)

	temp = strings.NewReader(errFile)
	_, err := ParseTemplate(temp)
	if err == nil || err.Error() != strErr {
		t.Error("Unsuccessive test-case is failed")
	}
}
