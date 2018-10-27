package et

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/crawlerclub/ce"
	"github.com/tkuchiki/parsetime"
	"golang.org/x/net/html"
	"zliu.org/goutil"
)

var timeParser, _ = parsetime.NewParseTime("Asia/Shanghai")

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
			item := root.Item["root"].(map[string]interface{})
			if len(item) > 0 {
				items = append(items, item)
			}
		}
	}

	if p.DefaultFields {
		for _, v := range items {
			v["from_url_"] = pageUrl
			v["from_parser_"] = p.Name
			v["crawl_time_"] = time.Now().UTC().Format(time.RFC3339)

			if v["time_"] != nil {
				switch v["time_"].(type) {
				case string:
					t, err := timeParser.Parse(v["time_"].(string))
					if err != nil {
						v["time_"] = t.UTC().Format(time.RFC3339)
					} else {
						delete(v, "time_")
					}
				case []interface{}:
					arr := v["time_"].([]interface{})
					if len(arr) == 0 {
						delete(v, "time_")
					} else {
						var t time.Time
						for _, ts := range arr {
							tt, err := timeParser.Parse(fmt.Sprintf("%+v", ts))
							if err == nil && tt.After(t) {
								t = tt
							}
						}
						if t.IsZero() {
							delete(v, "time_")
						} else {
							v["time_"] = t.UTC().Format(time.RFC3339)
						}
					}
				default:
					delete(v, "time_")
				}
			}
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
		if rule.Key == "" {
			return nil, nil, nil, fmt.Errorf("Key for %s is empty", name)
		}
		vals, err := p.parseNodeByRule(node, rule, pageUrl)
		if err != nil {
			return nil, nil, nil, err
		}
		if rule.Type == "dom" {
			for _, v := range vals {
				nodes = append(nodes, &DomNode{rule.Key, v.(*html.Node), item})
			}
		} else {
			if rule.Type == "url" {
				for _, v := range vals {
					switch v.(type) {
					case string:
						urls = append(urls,
							&UrlTask{ParserName: rule.Key, Url: v.(string)})
					case []map[string]interface{}:
						for _, m := range v.([]map[string]interface{}) {
							if m["url"] == nil || m["key"] == nil {
								continue
							}
							u, o1 := m["url"].(string)
							k, o2 := m["key"].(string)
							if o1 && o2 {
								urls = append(urls, &UrlTask{ParserName: k, Url: u})
							}
						}
					}
				}
			}

			if item[rule.Key] == nil {
				if len(vals) == 1 {
					item[rule.Key] = vals[0]
				} else if len(vals) > 1 {
					item[rule.Key] = interface{}(vals)
				}
			} else {
				switch item[rule.Key].(type) {
				case []interface{}:
					item[rule.Key] = append(
						item[rule.Key].([]interface{}), vals...)
				default:
					item[rule.Key] = append(
						[]interface{}{item[rule.Key]}, vals...)
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
		case "text":
			ret = append(ret, interface{}(ce.TextFromHTML(htmlquery.OutputHTML(n, true))))
		case "html":
			ret = append(ret, interface{}(htmlquery.OutputHTML(n, true)))
		case "attr":
			ret = append(ret, interface{}(htmlquery.SelectAttr(n, rule.Key)))
		default:
			return nil, fmt.Errorf("unknown rule type: %s", rule.Type)
		}
	}
	if rule.Re != nil {
		var vals []interface{}
		switch rule.Type {
		case "text", "html":
			for _, v := range ret {
				obj := make(map[string]string)
				for _, r := range rule.Re {
					m, err := goutil.RegexpExtract(v.(string), r)
					if err != nil {
						return nil, fmt.Errorf("regex:[%s] error:%+v", r, err)
					}
					for k, vv := range m {
						if rule.Type == "html" {
							obj[k] = html.UnescapeString(vv)
						} else {
							obj[k] = vv
						}
					}
				}
				switch len(obj) {
				case 0:
					vals = append(vals, v)
				case 1:
					for _, vv := range obj {
						vals = append(vals, interface{}(vv))
					}
				default:
					vals = append(vals, interface{}(obj))
				}
			}
			ret = vals
		case "url":
			for _, v := range ret {
				drop := false
				for _, r := range rule.Re {
					if !goutil.RegexpMatch(v.(string), r) {
						drop = true
						break
					}
				}
				if !drop {
					vals = append(vals, v)
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
