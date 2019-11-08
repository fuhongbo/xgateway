/**
* @Author: HongBo Fu
* @Date: 2019/10/23 08:21
 */

package node

import (
	"github.com/robertkrimen/otto"
	"net/http"
)

type Process interface {
	Exec(w http.ResponseWriter, r *http.Request, vm *otto.Otto) bool
}
