package et

import (
	"github.com/robertkrimen/otto"
	"golang.org/x/net/html"
	"sync"
)

type Rule struct {
	sync.Once
	// RuleTypes: url, dom, text, html
	Type    string `json:"type"`
	ItemKey string `json:"item_key"`
	Xpath   string `json:"xpath"`
	Re      string `json:"re"`
	Js      string `json:"js"`

	vm *otto.Otto
}

type Parser struct {
	sync.Once
	Name          string             `json:"name"`
	DefaultFields bool               `json:"default_fields""`
	ExampleUrl    string             `json:"example_url"`
	Rules         map[string][]*Rule `json:"rules"`
	Js            string             `json:"js"`

	vm *otto.Otto
}

type DomNode struct {
	Name string
	Node *html.Node
	Item map[string]interface{}
}

type UrlTask struct {
	ParserName string `json:"parser_name"`
	Url        string `json:"url"`
}
