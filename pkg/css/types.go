package css

import "github.com/deglebe/browse/pkg/dom"

type Stylesheet struct {
	Rules []Rule
}

type Rule struct {
	Selector	Selector
	Declarations	[]Declaration
}

type Declaration struct {
	Property	string
	Value		string
}

type Selector interface {
	Matches(el *dom.Node) bool
	Specificity() int
}
