/**
* @Author: HongBo Fu
* @Date: 2019/10/17 08:58
 */

package node

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type NodeServer struct {
	Server       *http.Server
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (node *NodeServer) Run() {
	port := ":" + node.Port

	//初始化HttpServer
	node.Server = &http.Server{
		Addr:         port,
		ReadTimeout:  node.ReadTimeout,
		WriteTimeout: node.WriteTimeout,
	}

	//注册Http Handler
	node.Server.Handler = http.HandlerFunc(Handle)

	logrus.WithFields(logrus.Fields{"Port": node.Port}).Info("Node Server Start...")

	err := node.Server.ListenAndServe()

	if err != nil {
		logrus.WithFields(logrus.Fields{"port": node.Port}).Error("Node Server 启动失败")
	}

}
