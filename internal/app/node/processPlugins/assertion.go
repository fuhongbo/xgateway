/**
* @Author: HongBo Fu
* @Date: 2019/10/29 15:43
 */

package processPlugins

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"net/http"
	"xgateway/internal/app/node/utility"
)

type assertionResult struct {
	result bool
	data   interface{}
	status int
}

type Assertion struct {
	Action string
	Value  string
	Type   string
}

func (p *Assertion) Exec(w http.ResponseWriter, r *http.Request, vm *otto.Otto) bool {
	if p.Type == "" {
		p.Type, p.Value = utility.SemanticDetection(p.Action)
	}
	switch p.Type {
	case "js":
		vm.Run(p.Value)
		v, _ := vm.Get("assertion")

		obj, _ := v.Object().Value().Export()
		result := obj.(map[string]interface{})

		if !result["result"].(bool) {
			w.WriteHeader(int(result["status"].(int64)))
			_, _ = fmt.Fprintln(w, result["data"])
			return false
		}

	}

	return true
}
