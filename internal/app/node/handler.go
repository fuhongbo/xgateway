/**
* @Author: HongBo Fu
* @Date: 2019/10/17 10:25
 */

package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"xgateway/internal/app/cfg"
	"xgateway/internal/app/node/utility"
	"xgateway/internal/pkg/syncmap"
)

var HandlerInstance syncmap.Map

func Handle(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI == "/favicon.ico" {
		return
	}

	if strings.Contains(r.URL.Path, "reload") {
		cfg.Reload("/Users/fuhongbo/go/src/xgateway/cmd/config.yaml", "yaml")
		ServiceOPS.StopHealth()
		ServiceOPS.Load()
		ServiceOPS.StartHealthCheck()
		ServiceOPS.RunLimiter()
		_, _ = fmt.Fprintln(w, "realod finish")
		return
	}

	if strings.Contains(r.URL.Path, "status") {
		data, err := json.Marshal(ServiceOPS.Services)

		if err != nil {
			fmt.Println(err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintln(w, string(data))
		return
	}

	baseURL := r.URL.Path[1:]
	baseURL = strings.Split(baseURL, "/")[0]

	//插件处理 -- 插件在配置文件中配置,启动时候进行初始化
	if !ProcessManager.Exec(baseURL, w, r) {
		return
	}

	if _, ok := ServiceOPS.Services[baseURL]; ok {

		t, v := utility.SemanticDetection(ServiceOPS.ServicesConfig[baseURL]["loadBalanceKey"].(string))

		switch t {

		case "bind":
			key, ok := utility.ReqeustParaser(v, r)

			if !ok {
				_, _ = fmt.Fprintln(w, v)
				return
			}
			server, err := ServiceOPS.Balance(baseURL, key)
			if err != nil {
				_, _ = fmt.Fprintln(w, err.Error())
				return
			}
			proxy(server, w, r)
		case "fix":

			server, err := ServiceOPS.Balance(baseURL, v)

			if err != nil {
				_, _ = fmt.Fprintln(w, err.Error())
				return
			}
			proxy(server, w, r)
		default:
			server, err := ServiceOPS.Balance(baseURL, "")
			if err != nil {
				_, _ = fmt.Fprintln(w, err.Error())
				return
			}
			proxy(server, w, r)
		}
	} else {
		_, _ = fmt.Fprintln(w, utility.Error_NotFound_Service)
		return
	}
}

func proxy(target string, w http.ResponseWriter, r *http.Request) {
	url, _ := url.Parse(target)

	var proxy *httputil.ReverseProxy
	if ins, ok := HandlerInstance.Load(url); ok {
		proxy = ins.(*httputil.ReverseProxy)
	} else {
		proxy = httputil.NewSingleHostReverseProxy(url)
		HandlerInstance.Store(url, proxy)
	}
	//proxy := httputil.NewSingleHostReverseProxy(url)

	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Header.Set("X-Service-ID", strings.Split(r.URL.Path, "/")[1])
	r.Host = url.Host

	r.URL.Path = "/" + strings.Join(strings.Split(r.URL.Path, "/")[2:], "/")
	//proxy.Transport = &customTransport{}
	proxy.ErrorHandler = ErrHandle

	//proxy.Director(r)
	proxy.ServeHTTP(w, r)

}

func ErrHandle(res http.ResponseWriter, req *http.Request, err error) {
	println(err.Error())
	res.WriteHeader(503)
	_, _ = fmt.Fprintln(res, "服务器暂时无法访问")
	return
}

type customTransport struct {
}

func (t *customTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := http.DefaultTransport.RoundTrip(request)

	//这里如果是发送的请求有错误，或者返回状态错误的话，应该是记录一个错误，同时应该有开关，是否记录这些错误会好一些

	if err != nil {

	}
	if response != nil {
		//println(request.Header.Get("X-Service-ID"))
	}

	return response, err
}
