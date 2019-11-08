/**
* @Author: HongBo Fu
* @Date: 2019/10/24 13:45
 */

package processPlugins

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"net/http"
	"strings"
	"xgateway/internal/app/node/ratelimit"
	"xgateway/internal/app/node/utility"
)

type RateLimiter struct {
	Calls      int
	TimePeriod int
	CounterKey string
	Type       string
	Value      string
}

func (p *RateLimiter) Exec(w http.ResponseWriter, r *http.Request, vm *otto.Otto) bool {
	if p.Type == "" {
		p.Type, p.Value = utility.SemanticDetection(p.CounterKey)
	}
	counterKey := ""
	switch p.Type {
	case "js":
		vm.Run(p.Value)
		v, _ := vm.Get("counterKey")
		counterKey = v.String()
	case "bind":
		v, ok := utility.ReqeustParaser(p.Value, r)

		if ok {
			counterKey = v
		} else {
			_, _ = fmt.Fprintln(w, v)
			return false
		}

	case "fix":
		counterKey = p.Value
	}

	baseURL := r.URL.Path[1:]
	baseURL = strings.Split(baseURL, "/")[0]

	if !ratelimit.LimiterInstance.GetLimiter().TokenAccess(counterKey, baseURL) {
		w.WriteHeader(429)
		_, _ = fmt.Fprintln(w, utility.Error_Rate_Limit)
		return false
	}
	return true
}
