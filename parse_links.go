package et

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"zliu.org/goutil"
)

func ParseLinks(page, url string) ([]string, error) {
	doc, err := htmlquery.Parse(strings.NewReader(page))
	if err != nil {
		return nil, fmt.Errorf("htmlquery.Parse err: %+v", err)
	}
	var links []string
	htmlquery.FindEach(doc, "//a", func(i int, node *html.Node) {
		link := htmlquery.SelectAttr(node, "href")
		if u, err := goutil.MakeAbsoluteUrl(link, url); err == nil {
			if strings.HasPrefix(u, "http") && !strings.HasSuffix(u, ".exe") {
				links = append(links, u)
			}
		}
	})
	return links, nil
}
