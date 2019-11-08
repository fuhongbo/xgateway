/**
* @Author: HongBo Fu
* @Date: 2019/10/23 14:11
 */

package processPlugins

import (
	"github.com/robertkrimen/otto"
	"net/http"
	"xgateway/internal/app/node/utility"
)

type VariableSet struct {
	Name  string
	Value string
	Type  string
}

func (p *VariableSet) Exec(w http.ResponseWriter, r *http.Request, vm *otto.Otto) bool {
	if p.Type == "" {
		p.Type, p.Value = utility.SemanticDetection(p.Value)
	}
	switch p.Type {
	case "js":
		vm.Run(p.Value)
	}

	return true
}
