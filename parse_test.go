package et

import (
	"crawler.club/dl"
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestParse(t *testing.T) {
	b, err := ioutil.ReadFile("./newsmth.net.json")
	if err != nil {
		t.Fatal(err)
	}
	p := new(Parser)
	if err = json.Unmarshal(b, p); err != nil {
		t.Fatal(err)
	}
	url := "http://www.newsmth.net/nForum/article/Taiwan/50328"
	resp := dl.DownloadUrl(url)
	if resp.Error != nil {
		t.Fatal(resp.Error)
	}
	urls, items, err := p.Parse(resp.Text, url)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(urls)
	t.Log(items)
}
