package et

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

func (r *Rule) RunJs(v interface{}) (interface{}, error) {
	if r.Js == "" {
		return v, nil
	}
	var err error
	r.Do(func() {
		r.vm = otto.New()
		if _, err = r.vm.Run(r.Js); err != nil {
			r.vm = nil
			return
		}
	})
	if err != nil {
		return nil, err
	}
	if r.vm == nil {
		return nil, fmt.Errorf("rule.vm is nil")
	}
	jsVal, err := r.vm.ToValue(v)
	if err != nil {
		return nil, err
	}
	ret, err := r.vm.Call("process", nil, jsVal)
	if err != nil {
		return nil, err
	}
	return ret.Export()
}
