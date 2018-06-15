package et

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"crawler.club/dl"
	"zliu.org/goutil"
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
	url := "http://www.newsmth.net/nForum/article/RealEstate/7124059"
	resp := dl.DownloadUrl(url)
	if resp.Error != nil {
		t.Fatal(resp.Error)
	}
	urls, items, err := p.Parse(resp.Text, url)
	if err != nil {
		t.Fatal(err)
	}
	if b, err = goutil.JsonMarshalIndent(urls, "", "  "); err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
	if b, err = goutil.JsonMarshalIndent(items, "", "  "); err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
