package et

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"net/url"
	"strings"
)

func (p *Parser) Parse(page, pageUrl string) error {
	if page == "" {
		return fmt.Errorf("page content is empty")
	}
	u, err := url.Parse(pageUrl)
	if err != nil || !u.IsAbs() {
		return fmt.Errorf("pageUrl %s is not abs url", pageUrl)
	}
	doc, err := htmlquery.Parse(strings.NewReader(page))
	if err != nil {
		return fmt.Errorf("htmlquery.Parse err: %+v", err)
	}

	root := &DomNode{"root", doc, make(map[string]interface{})}
	domList := []*DomNode{root}
	for {
		if len(domList) == 0 {
			break
		}
		name := domList[0].Name
		node := domList[0].Node
		//parents := domList[0].Item
		domList = domList[1:]
		//nodes, urls, item, err :=
		p.parseNode(node, name, pageUrl)
	}
	return nil
}

func (p *Parser) parseNode(node *html.Node, name string, pageUrl string) error {
	if node == nil {
		return fmt.Errorf("node is nil in parseNode")
	}
	if p.Rules[name] == nil {
		return fmt.Errorf("parse rule for %s not found", name)
	}
	for _, rule := range p.Rules[name] {
		if rule.ItemKey == "" {
			return fmt.Errorf("ItemKey for %s is empty", name)
		}
		p.parseNodeByRule(node, rule, pageUrl)
	}
	return nil
}

func (p *Parser) parseNodeByRule(node *html.Node, rule *Rule, pageUrl string) {
}
