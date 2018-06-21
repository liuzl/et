package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"crawler.club/dl"
	"crawler.club/et"
	"zliu.org/goutil"
)

var (
	conf = flag.String("conf", "bjjtgl.json", "parse conf file")
	url  = flag.String("url", "http://sslk.bjjtgl.gov.cn/roadpublish/Map/vmsimg/vmsimage/secondrealinfo.htm", "url")
)

func main() {
	flag.Parse()
	b, err := ioutil.ReadFile(*conf)
	if err != nil {
		log.Fatal(err)
	}
	p := new(et.Parser)
	if err = json.Unmarshal(b, p); err != nil {
		log.Fatal(err)
	}
	resp := dl.DownloadUrl(*url)
	if resp.Error != nil {
		log.Fatal(resp.Error)
	}
	urls, items, err := p.Parse(resp.Text, *url)
	if err != nil {
		log.Fatal(err)
	}
	if b, err = goutil.JsonMarshalIndent(urls, "", "  "); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tasks: %s\n", string(b))
	if b, err = goutil.JsonMarshalIndent(items, "", "  "); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("items: %s\n", string(b))
}
