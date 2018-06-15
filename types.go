package et

import (
	"golang.org/x/net/html"
)

type Rule struct {
	// RuleTypes: url, dom, text, html
	Type  string   `json:"type"`
	Key   string   `json:"key"`
	Xpath string   `json:"xpath"`
	Re    []string `json:"re"`
	Js    string   `json:"js"`
}

type Parser struct {
	Name          string             `json:"name"`
	DefaultFields bool               `json:"default_fields""`
	ExampleUrl    string             `json:"example_url"`
	Rules         map[string][]*Rule `json:"rules"`
	Js            string             `json:"js"`
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
