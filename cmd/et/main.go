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
	url  = flag.String("url", "", "url to download and parse")
	ua   = flag.String("ua", "", "pc, mobile, google. Golang UA for empty")
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
	var u = p.ExampleUrl
	if *url != "" {
		u = *url
	}
	var req = &dl.HttpRequest{Url: u}
	if *ua == "pc" || *ua == "mobile" || *ua == "google" {
		req.Platform = *ua
	}
	resp := dl.Download(req)
	if resp.Error != nil {
		log.Fatal(resp.Error)
	}

	urls, items, err := p.Parse(resp.Text, u)
	if err != nil {
		log.Fatal(err)
	}
	if b, err = goutil.JSONMarshalIndent(urls, "", "  "); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tasks: %s\n", string(b))
	if b, err = goutil.JSONMarshalIndent(items, "", "  "); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("items: %s\n", string(b))
}
