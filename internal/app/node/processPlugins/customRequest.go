/**
* @Author: HongBo Fu
* @Date: 2019/10/28 08:30
 */

package processPlugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"net/http"
	"time"
	"xgateway/internal/app/node/utility"
)

type CustomRequest struct {
	Identification string
	URL            string
	Method         string
	Header         map[string]string
	Body           string
	Timeout        int
}

func (p *CustomRequest) Exec(w http.ResponseWriter, r *http.Request, vm *otto.Otto) bool {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	bd, ok := vmExec("body", p.Body, w, r, vm)

	if bd == "undefined" {
		w.WriteHeader(503)
		_, _ = fmt.Fprintln(w, "自定义认证请求Body配置错误")
		return false
	}

	if !ok {
		w.WriteHeader(503)
		_, _ = fmt.Fprintln(w, "自定义认证请求Body配置错误")
		return false
	}
	req, _ := http.NewRequest(p.Method, p.URL, bytes.NewReader([]byte(bd)))
	client.Timeout = time.Duration(p.Timeout) * time.Second
	for k, v := range p.Header {
		result, ok := vmExec(k, v, w, r, vm)
		if !ok {
			w.WriteHeader(503)
			_, _ = fmt.Fprintln(w, "自定义认证请求头部请求配置错误")
			return false
		}
		req.Header.Add(k, result)
	}

	resp, err := client.Do(req)

	if err != nil {
		w.WriteHeader(503)
		_, _ = fmt.Fprintln(w, "自定义请求访问发生错误 - ", err.Error())
		return false
	}

	defer resp.Body.Close()

	rbody, _ := ioutil.ReadAll(resp.Body)

	tempResp := make(map[string]interface{})
	tempResp["status"] = resp.StatusCode
	tempResp["body"] = string(rbody)

	result := struct {
		Response map[string]interface{}
	}{
		tempResp,
	}

	vm.Set(p.Identification, result)

	return true
}

func vmExec(key, script string, w http.ResponseWriter, r *http.Request, vm *otto.Otto) (string, bool) {

	rtype, value := utility.SemanticDetection(script)

	result := ""
	switch rtype {
	case "js":

		vm.Run(value)
		v, _ := vm.Get(key)
		if v.IsObject() {
			obj, _ := v.Object().Value().Export()
			jsondata, _ := json.Marshal(obj)
			result = string(jsondata)
		} else {
			result = v.String()
		}
	case "bind":
		v, ok := utility.ReqeustParaser(value, r)

		if ok {
			result = v
		} else {
			return "", false
		}

	case "fix":
		result = value
	}

	return result, true
}
