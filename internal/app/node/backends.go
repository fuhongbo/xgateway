/**
* @Author: HongBo Fu
* @Date: 2019/10/17 13:42
 */

package node

import (
	"net/http"
	"time"
	"xgateway/internal/app/node/ratelimit"
)

type BackEnds struct {
	//服务器地址
	EndPoint string `json:"endPoint"`
	//健康状态
	Health bool `json:"health"`
	//健康检查URL
	CheckURL string `json:"checkURL"`
	//健康检查间隔
	CheckTicker int `json:"checkTicker"`
	//是否结束健康检查
	Done chan int `json:"-"`
}

func (b *BackEnds) CheckHealth(serviceName string) {
	return
	if b.Health {
		tick := time.NewTicker(time.Duration(b.CheckTicker) * time.Second)
		for {
			select {
			case <-b.Done:
				tick.Stop()
			case <-tick.C:
				//logrus.Info(b.EndPoint, "  ---> 健康检查中...")

				go func(temp *BackEnds) {
					healthURL := b.EndPoint + b.CheckURL
					client := &http.Client{}
					defer client.CloseIdleConnections()
					reqest, err := http.NewRequest("GET", healthURL, nil)
					if err != nil {
						if !ratelimit.LimiterInstance.GetLimiter().TokenAccess(b.EndPoint, serviceName+"_health") {
							ServiceOPS.SetHealthStatus(serviceName, false, b.EndPoint)
						}
						return
					}
					client.Timeout = 10 * time.Second
					response, err := client.Do(reqest)
					if err != nil {
						if !ratelimit.LimiterInstance.GetLimiter().TokenAccess(b.EndPoint, serviceName+"_health") {
							ServiceOPS.SetHealthStatus(serviceName, false, b.EndPoint)
						}
						return
					}
					if response.StatusCode != 200 {
						if !ratelimit.LimiterInstance.GetLimiter().TokenAccess(b.EndPoint, serviceName+"_health") {
							ServiceOPS.SetHealthStatus(serviceName, false, b.EndPoint)
						}
						return
					}

					ServiceOPS.SetHealthStatus(serviceName, true, b.EndPoint)
				}(b)
			}
		}
	}
}
