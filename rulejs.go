package et

import (
	"github.com/robertkrimen/otto"
)

// RunJs runs the rule's js code for v
func (r *Rule) RunJs(v interface{}) (interface{}, error) {
	if r.Js == "" {
		return v, nil
	}
	var err error
	vm := otto.New()
	if _, err = vm.Run(r.Js); err != nil {
		return nil, err
	}
	jsVal, err := vm.ToValue(v)
	if err != nil {
		return nil, err
	}
	ret, err := vm.Call("process", nil, jsVal)
	if err != nil {
		return nil, err
	}
	return ret.Export()
}
