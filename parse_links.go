package et

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/liuzl/store"
	"golang.org/x/net/html"
	"zliu.org/goutil"
)

var linkStore *store.LevelStore
var once sync.Once

func getLinkStore() *store.LevelStore {
	once.Do(func() {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			panic(err)
		}
		linkStore, err = store.NewLevelStore(filepath.Join(dir, ".etlinks"))
		if err != nil {
			panic(err)
		}
	})
	return linkStore
}

// ParseNewLinks returns new urls contained in html page
func ParseNewLinks(page, url string) ([]string, error) {
	links, err := ParseLinks(page, url)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, link := range links {
		has, err := getLinkStore().Has(link)
		if err != nil {
			return nil, err
		}
		if has {
			continue
		}
		getLinkStore().Put(link, []byte(time.Now().UTC().Format(time.RFC3339)))
		ret = append(ret, link)
	}
	return ret, nil
}

// ParseLinks returns all urls contained in html page
func ParseLinks(page, url string) ([]string, error) {
	doc, err := htmlquery.Parse(strings.NewReader(page))
	if err != nil {
		return nil, fmt.Errorf("htmlquery.Parse err: %+v", err)
	}
	var links []string
	htmlquery.FindEach(doc, "//a", func(i int, node *html.Node) {
		link := htmlquery.SelectAttr(node, "href")
		if u, err := goutil.MakeAbsoluteURL(link, url); err == nil {
			if strings.HasPrefix(u, "http") && !strings.HasSuffix(u, ".exe") {
				links = append(links, u)
			}
		}
	})
	return links, nil
}
