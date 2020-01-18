package et

import (
	"crawler.club/dl"
)

func (p *Parser) ParseURL(url string) ([]*UrlTask, []map[string]interface{}, error) {
	req := &dl.HttpRequest{Url: url}
	resp := dl.Download(req)
	if resp.Error != nil {
		return nil, nil, resp.Error
	}
	return p.Parse(resp.Text, url)
}
