package et

import (
	"golang.org/x/net/html"
)

// Rule extract a specific key by xpath, regexp and js sequentially.
// Five types for now: url, dom, text, html and attr
type Rule struct {
	Type  string   `json:"type"`
	Key   string   `json:"key"`
	Xpath string   `json:"xpath"`
	Re    []string `json:"re"`
	Js    string   `json:"js"`
}

// Parser contains a set of cascaded rule and an optional js code to parse
// corresponding htmls
type Parser struct {
	Name          string             `json:"name"`
	DefaultFields bool               `json:"default_fields"`
	ExampleUrl    string             `json:"example_url"`
	UA            string             `json:"ua"`
	Urls          []string           `json:"urls"`
	Rules         map[string][]*Rule `json:"rules"`
	Js            string             `json:"js"`
}

// DomNode is for internal usage
type DomNode struct {
	Name string
	Node *html.Node
	Item map[string]interface{}
}

// UrlTask contains a crawling task of Url that should be parsed by ParserName
type UrlTask struct {
	ParserName string      `json:"parser_name"`
	Url        string      `json:"url"`
	TaskName   string      `json:"task_name"`
	Ext        interface{} `json:"ext"`
}
