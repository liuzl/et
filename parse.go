package et

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/liuzl/goutil"
	"golang.org/x/net/html"
	"net/url"
	"strings"
	"time"
)

func (p *Parser) Parse(
	page, pageUrl string) ([]*UrlTask, []map[string]interface{}, error) {

	if page == "" {
		return nil, nil, fmt.Errorf("page content is empty")
	}
	u, err := url.Parse(pageUrl)
	if err != nil || !u.IsAbs() {
		return nil, nil, fmt.Errorf("pageUrl %s is not abs url", pageUrl)
	}
	doc, err := htmlquery.Parse(strings.NewReader(page))
	if err != nil {
		return nil, nil, fmt.Errorf("htmlquery.Parse err: %+v", err)
	}

	root := &DomNode{"root", doc, make(map[string]interface{})}
	domList := []*DomNode{root}
	var items []map[string]interface{}
	var urlList []*UrlTask
	for {
		if len(domList) == 0 {
			break
		}
		name := domList[0].Name
		node := domList[0].Node
		parents := domList[0].Item
		domList = domList[1:]
		nodes, urls, item, err := p.parseNode(node, name, pageUrl)
		if err != nil {
			return nil, nil, err
		}
		domList = append(domList, nodes...)
		urlList = append(urlList, urls...)
		if item != nil {
			if parents[name] == nil {
				parents[name] = interface{}(item)
			} else {
				switch parents[name].(type) {
				case []interface{}:
					parents[name] = append(parents[name].([]interface{}),
						interface{}(item))
				default:
					parents[name] = []interface{}{parents[name], interface{}(item)}
				}
				parents[name] = interface{}(parents[name])
			}
		}
	} // end for

	if root.Item["root"] != nil {
		switch root.Item["root"].(type) {
		case []interface{}:
			t, _ := root.Item["root"].([]map[string]interface{})
			items = append(items, t...)
		default:
			items = append(items, root.Item["root"].(map[string]interface{}))
		}
	}

	if p.DefaultFields {
		for _, v := range items {
			v["from_url_"] = pageUrl
			v["from_parser_"] = p.Name
			v["crawl_time_"] = time.Now().Format("2006-01-02 15:04:05")
		}
	}
	items, err = p.RunJs(items)
	return urlList, items, err
}

func (p *Parser) parseNode(node *html.Node, name string,
	pageUrl string) ([]*DomNode, []*UrlTask, map[string]interface{}, error) {

	if node == nil {
		return nil, nil, nil, fmt.Errorf("node is nil in parseNode")
	}
	if p.Rules[name] == nil {
		return nil, nil, nil, fmt.Errorf("parse rule for %s not found", name)
	}

	var nodes []*DomNode
	var urls []*UrlTask
	item := make(map[string]interface{})

	for _, rule := range p.Rules[name] {
		if rule.ItemKey == "" {
			return nil, nil, nil, fmt.Errorf("ItemKey for %s is empty", name)
		}
		vals, err := p.parseNodeByRule(node, rule, pageUrl)
		if err != nil {
			return nil, nil, nil, err
		}
		if rule.Type == "dom" {
			for _, v := range vals {
				nodes = append(nodes, &DomNode{rule.ItemKey, v.(*html.Node), item})
			}
		} else {
			if rule.Type == "url" {
				for _, v := range vals {
					if u, ok := v.(string); ok {
						urls = append(urls, &UrlTask{rule.ItemKey, u})
					}
				}
			}

			if item[rule.ItemKey] == nil {
				if len(vals) == 1 {
					item[rule.ItemKey] = vals[0]
				} else if len(vals) > 1 {
					item[rule.ItemKey] = interface{}(vals)
				}
			} else {
				switch item[rule.ItemKey].(type) {
				case []interface{}:
					item[rule.ItemKey] = append(
						item[rule.ItemKey].([]interface{}), vals...)
				default:
					item[rule.ItemKey] = append(
						[]interface{}{item[rule.ItemKey]}, vals...)
				}
			}
		}
	}
	return nodes, urls, item, nil
}

func (p *Parser) parseNodeByRule(
	node *html.Node, rule *Rule, pageUrl string) ([]interface{}, error) {

	if node == nil {
		return nil, fmt.Errorf("node is nil")
	}
	if rule == nil || rule.Type == "" || rule.Xpath == "" {
		return nil, fmt.Errorf("invalid rule: %+v", rule)
	}

	var ret []interface{}
	for _, n := range htmlquery.Find(node, rule.Xpath) {
		switch rule.Type {
		case "dom":
			ret = append(ret, interface{}(n))
		case "url":
			if u, err := goutil.MakeAbsoluteUrl(
				htmlquery.SelectAttr(n, "href"), pageUrl); err != nil {
				return nil, fmt.Errorf("MakeAbsoluteUrl err: %+v", err)
			} else {
				ret = append(ret, interface{}(u))
			}
		case "string":
			ret = append(ret, interface{}(htmlquery.InnerText(n)))
		case "html":
			ret = append(ret, interface{}(htmlquery.OutputHTML(n, true)))
		}
	}
	if rule.Re != "" {
		var vals []interface{}
		switch rule.Type {
		case "string":
			for _, v := range ret {
				res, err := goutil.RegexpParse(v.(string), rule.Re)
				if err != nil {
					return nil, fmt.Errorf("Re:[%s] error: %+v", rule.Re, err)
				}
				for _, i := range res {
					vals = append(vals, interface{}(i))
				}
			}
			ret = vals
		case "url":
			for _, v := range ret {
				if goutil.RegexpMatch(v.(string), rule.Re) {
					vals = append(vals, interface{}(v))
				}
			}
			ret = vals
		}
	}

	if rule.Js != "" {
		var vals []interface{}
		for _, v := range ret {
			s, err := rule.RunJs(v)
			if err != nil {
				return nil, fmt.Errorf("rule.RuleJs error: %+v", err)
			}
			vals = append(vals, s)
		}
		ret = vals
	}

	return ret, nil
}
