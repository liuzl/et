package et

import (
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"crawler.club/dl"
)

var (
	storeUri = flag.String("store_uri", "127.0.0.1:2002", "parser store addr")
)

type Parsers struct {
	sync.Mutex
	items map[string]*Parser
}

func (p *Parsers) GetParser(name string) (*Parser, error) {
	p.Lock()
	defer p.Unlock()
	if p.items[name] != nil {
		return p.items[name], nil
	}
	resp := dl.DownloadUrl(fmt.Sprintf("http://%s/get/%s", *storeUri, name))
	if resp.Error != nil {
		return nil, resp.Error
	}
	parser := new(Parser)
	if err := json.Unmarshal(resp.Content, parser); err != nil {
		return nil, err
	}
	p.items[name] = parser
	return parser, nil
}

var pool = &Parsers{items: make(map[string]*Parser)}

func Parse(name, page, url string) ([]*UrlTask, []map[string]interface{}, error) {
	p, err := pool.GetParser(name)
	if err != nil {
		return nil, nil, err
	}
	return p.Parse(page, url)
}
