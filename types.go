package et

import (
	"golang.org/x/net/html"
)

type Rule struct {
	// four RuleTypes: url, dom, string, html
	Type    string `json:"type"`
	ItemKey string `json:"item_key"`
	Xpath   string `json:"xpath"`
	Regex   string `json:"regex"`
	Js      string `json:"js"`
}

type Parser struct {
	Name          string             `json:"name"`
	DefaultFields bool               `json:"default_fields""`
	ExampleUrl    string             `json:"example_url"`
	Rules         map[string][]*Rule `json:"rules"`
	PostProcessor string             `json:"post_processor"`
}

type DomNode struct {
	Name string
	Node *html.Node
	Item map[string]interface{}
}
