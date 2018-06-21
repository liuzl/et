package et

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"crawler.club/dl"
	"zliu.org/goutil"
)

var tests = [][]string{
	//[]string{"article.json", "http://www.newsmth.net/nForum/article/RealEstate/7124059"},
	//[]string{"section.json", "http://www.newsmth.net/nForum/section/1"},
	[]string{"board.json", "http://www.newsmth.net/nForum/board/Universal"},
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		b, err := ioutil.ReadFile(test[0])
		if err != nil {
			t.Fatal(err)
		}
		p := new(Parser)
		if err = json.Unmarshal(b, p); err != nil {
			t.Fatal(err)
		}
		url := test[1]
		resp := dl.DownloadUrlWithProxy(url)
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
}
