package et

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

func (p *Parser) RunJs(
	items []map[string]interface{}) ([]map[string]interface{}, error) {

	if p.Js == "" {
		return items, nil
	}
	var err error
	p.Do(func() {
		p.vm = otto.New()
		if _, err = p.vm.Run(p.Js); err != nil {
			p.vm = nil
			return
		}
	})
	if err != nil {
		return nil, err
	}
	if p.vm == nil {
		return nil, fmt.Errorf("parser.vm is nil")
	}
	jsVal, err := p.vm.ToValue(items)
	if err != nil {
		return nil, err
	}
	ret, err := p.vm.Call("process", nil, jsVal)
	if err != nil {
		return nil, err
	}
	s, err := ret.Export()
	if err != nil {
		return nil, err
	}
	if value, ok := s.([]map[string]interface{}); ok {
		return value, nil
	}
	return nil, fmt.Errorf("s.([]map[string]interface{}) error")
}
