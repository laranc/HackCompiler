package compiler

import (
	"container/list"
)

var (
	tokens *list.List
)

type ParseTree struct {
	Type, Value string
	Children    *list.List
}

func NewParseTree(t, v string) *ParseTree {
	return &ParseTree{
		Type:     t,
		Value:    v,
		Children: list.New(),
	}
}

func (tree *ParseTree) addChild(t, v string) {
	tree.Children.PushBack(NewParseTree(t, v))
}

type ParseException struct{}

func (e *ParseException) Error() string {
	return "An error occurred while parsing"
}

func beginParsing(tokenList *list.List) {
	tokens = tokenList
	parseProgram()
	tokens.Init()
}

func parseProgram() {
}
