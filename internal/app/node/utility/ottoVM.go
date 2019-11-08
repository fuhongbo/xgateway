/**
* @Author: HongBo Fu
* @Date: 2019/10/23 15:04
 */

package utility

import (
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"strings"
)

var OTTOVM *VM

type VM struct {
}

func init() {
	OTTOVM = &VM{}
}

func (v *VM) GetDefaultVM(w http.ResponseWriter, r *http.Request) *otto.Otto {
	vm := otto.New()

	vm.Set("requestGet", func(call otto.FunctionCall) otto.Value {
		key, _ := call.Argument(0).ToString()

		v, ok := ReqeustParaser(key, r)
		if ok {
			result, _ := vm.ToValue(v)
			return result
		}

		return otto.Value{}
	})

	vm.Set("headerSet", func(call otto.FunctionCall) otto.Value {
		key, _ := call.Argument(0).ToString()
		value, _ := call.Argument(1).ToString()

		r.Header.Set(key, value)

		return otto.Value{}
	})

	vm.Set("bodySet", func(call otto.FunctionCall) otto.Value {
		new_body_content, _ := call.Argument(0).ToString()
		r.ContentLength = int64(len(new_body_content))
		r.Body = ioutil.NopCloser(strings.NewReader(new_body_content))
		return otto.Value{}
	})

	return vm
}
