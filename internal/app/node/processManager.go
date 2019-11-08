/**
* @Author: HongBo Fu
* @Date: 2019/10/23 08:28
 */

package node

import (
	"github.com/SongLiangChen/RateLimiter"
	"github.com/robertkrimen/otto"
	"net/http"
	"xgateway/internal/app/cfg"
	"xgateway/internal/app/node/processPlugins"
	"xgateway/internal/app/node/ratelimit"
	"xgateway/internal/app/node/utility"
)

type processManagerData struct {
	processMP map[string][]Process
}

var ProcessManager *processManagerData

func init() {
	ProcessManager = &processManagerData{}
	ProcessManager.processMP = make(map[string][]Process)
}

//添加新的流程处理
func (p *processManagerData) add(serviceName string, process Process) {

	p.processMP[serviceName] = append(p.processMP[serviceName], process)
}

func (p *processManagerData) Register(serviceName string, plugins []cfg.Plugin) {
	for _, v := range plugins {
		switch v.ActionType {
		case "authBasic":
			plugin := &processPlugins.BasicAuth{}
			plugin.UserName = v.Properties["userName"].(string)
			plugin.Password = v.Properties["password"].(string)
			p.add(serviceName, plugin)
		case "setVariable":
			plugin := &processPlugins.VariableSet{}
			plugin.Name = v.Properties["name"].(string)
			plugin.Value = v.Properties["value"].(string)
			p.add(serviceName, plugin)
		case "request":
			plugin := &processPlugins.CustomRequest{}
			plugin.URL = v.Properties["url"].(string)
			plugin.Method = v.Properties["method"].(string)
			plugin.Body = v.Properties["body"].(string)
			plugin.Timeout = v.Properties["timeout"].(int)
			plugin.Header = make(map[string]string)
			plugin.Identification = v.Properties["identification"].(string)
			for k, v := range v.Properties["header"].(map[interface{}]interface{}) {
				plugin.Header[k.(string)] = v.(string)
			}
			p.add(serviceName, plugin)
		case "rateLimit":
			plugin := &processPlugins.RateLimiter{}
			plugin.Calls = v.Properties["calls"].(int)
			plugin.TimePeriod = v.Properties["timePeriod"].(int)
			plugin.CounterKey = v.Properties["counterKey"].(string)
			ratelimit.LimiterInstance.AddRule(serviceName, &RateLimiter.Rule{
				Duration: int32(plugin.TimePeriod),
				Limit:    int32(plugin.Calls),
			})
			p.add(serviceName, plugin)
		case "ipFilter":
			plugin := &processPlugins.IPFilter{}
			plugin.Action = v.Properties["action"].(string)
			plugin.Address = v.Properties["address"].(string)
			p.add(serviceName, plugin)
		case "assertion":
			plugin := &processPlugins.Assertion{}
			plugin.Action = v.Properties["action"].(string)
			p.add(serviceName, plugin)
		case "jwt":
			plugin := &processPlugins.JwtAuth{}
			plugin.Type = v.Properties["type"].(string)
			plugin.Method = v.Properties["method"].(string)
			plugin.Kid = v.Properties["kid"].(string)
			plugin.PublicKey = v.Properties["publicKey"].(map[interface{}]interface{})
			plugin.Verify = v.Properties["verify"].(map[interface{}]interface{})
			p.add(serviceName, plugin)
		}
	}
}

func (p *processManagerData) Exec(serviceName string, w http.ResponseWriter, r *http.Request) bool {

	var vm *otto.Otto
	if len(p.processMP) > 0 {
		vm = utility.OTTOVM.GetDefaultVM(w, r)
	}
	for _, item := range p.processMP[serviceName] {
		if !item.Exec(w, r, vm) {
			return false
		}
	}
	return true

}
