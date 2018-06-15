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
	vm := otto.New()
	if _, err = vm.Run(p.Js); err != nil {
		return nil, err
	}
	jsVal, err := vm.ToValue(items)
	if err != nil {
		return nil, err
	}
	ret, err := vm.Call("process", nil, jsVal)
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
