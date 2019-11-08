/**
* @Author: HongBo Fu
* @Date: 2019/10/17 11:47
 */

package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"time"
	"xgateway/internal/app/cfg"
	"xgateway/internal/app/node"
)

func main() {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	//runtime.GOMAXPROCS(1)

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	var configFile = "/Users/fuhongbo/go/src/xgateway/cmd/config.yaml"

	//加载配置文件
	cfg.Reload(configFile, "yaml")
	node.ServiceOPS.Load()
	node.ServiceOPS.StartHealthCheck()

	//测试插件
	for _, item := range cfg.Instance.Config.Routes {
		node.ProcessManager.Register(strings.Trim(item.Route, "/"), item.ProcessPlugin.InBound)
	}

	node.ServiceOPS.RunLimiter()
	node.ServiceOPS.InitBalance()
	server := node.NodeServer{
		Server:       &http.Server{},
		Port:         cfg.Instance.Config.Port,
		ReadTimeout:  time.Duration(cfg.Instance.Config.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Instance.Config.WriteTimeout) * time.Second,
	}

	go func() {
		server.Run()
	}()

	// 限制 CPU 使用数，避免过载
	runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
	http.ListenAndServe("0.0.0.0:6060", nil)

}
