/**
* @Author: HongBo Fu
* @Date: 2019/10/17 14:23
 */

package node

import (
	"github.com/SongLiangChen/RateLimiter"
	"strings"
	"sync"
	"xgateway/internal/app/cfg"
	"xgateway/internal/app/node/ratelimit"
)

var ServiceOPS *Service

type Service struct {
	Services        map[string][]BackEnds
	ServicesConfig  map[string]map[string]interface{}
	BalanceInstance map[string]Balance
	sync.Mutex
}

func init() {
	ServiceOPS = &Service{}
	ServiceOPS.Services = make(map[string][]BackEnds)
	ServiceOPS.ServicesConfig = make(map[string]map[string]interface{})
	ServiceOPS.BalanceInstance = make(map[string]Balance)
}

func (s *Service) SetHealthStatus(serviceName string, status bool, endPoint string) {
	s.Lock()
	for i, item := range s.Services[serviceName] {
		if item.EndPoint == endPoint {
			s.Services[serviceName][i].Health = status
			if !status {
				s.RemoveEndPointInBalance(serviceName, endPoint)
			} else {
				s.AddEndPointInBalance(serviceName, endPoint)
			}
		}
	}
	s.Unlock()
}

func (s *Service) Load() {
	ServiceOPS.Lock()
	defer ServiceOPS.Unlock()

	for _, item := range cfg.Instance.Config.Routes {
		backends := []BackEnds{}
		serviceName := strings.Trim(item.Route, "/")
		for _, endpoint := range item.Endpoints {
			ep := BackEnds{}
			ep.EndPoint = endpoint
			ep.CheckTicker = item.CheckTicket
			ep.CheckURL = item.CheckHealthURL
			ep.Health = true
			ep.Done = make(chan int)
			backends = append(backends, ep)
		}
		ServiceOPS.AddOrUpdate(serviceName, backends)
		config := make(map[string]interface{})
		config["loadBalanceMode"] = item.NLBModel
		config["loadBalanceKey"] = item.LoadBalanceKey
		config["swagger"] = item.Swagger
		//config["authPlugin"] = item.AuthPlugin
		config["healthLimitCount"] = item.HealthLimitCount
		config["healthLimitTimeWindow"] = item.HealthLimitTimeWindow
		config["requestLimitCount"] = item.RequestLimitCount
		config["requestLimitTimeWindow"] = item.RequestLimitTimeWindow

		ServiceOPS.ServicesConfig[serviceName] = config
	}
}

func (s *Service) IsExist(serviceName string) bool {
	s.Lock()
	defer s.Unlock()
	if len(s.Services[serviceName]) > 0 {
		return true
	} else {
		return false
	}
}

func (s *Service) RunLimiter() {

	for k, _ := range s.Services {
		ratelimit.LimiterInstance.AddRule(k+"_health", &RateLimiter.Rule{
			Duration: s.ServicesConfig[k]["healthLimitTimeWindow"].(int32),
			Limit:    s.ServicesConfig[k]["healthLimitCount"].(int32),
		})
	}

	ratelimit.LimiterInstance.Start()

}

func (s *Service) StartHealthCheck() {
	for k, v := range s.Services {
		for _, item := range v {
			//go item.CheckHealth()
			go func(ser BackEnds, serviceName string) {
				ser.CheckHealth(serviceName)
			}(item, k)
		}
	}
}
func (s *Service) StopHealth() {
	for k, _ := range s.Services {
		s.StopHealthCheck(k)
	}
}
func (s *Service) StopHealthCheck(serviceName string) {
	for _, item := range s.Services[serviceName] {
		if item.Health {
			item.Done <- 1
		}
	}
}

func (s *Service) Remove(serviceName string) {
	s.StopHealthCheck(serviceName)
	delete(s.Services, serviceName)
}

func (s *Service) AddOrUpdate(serviceName string, backEnds []BackEnds) {
	s.Services[serviceName] = backEnds
}

func (s *Service) InitBalance() {

	for key, item := range s.ServicesConfig {
		bl := Balance{}
		bl.New(item["loadBalanceMode"].(string))
		for _, s := range s.Services[key] {
			bl.Add(s.EndPoint)
		}
		s.BalanceInstance[key] = bl
	}
}

func (s *Service) RemoveEndPointInBalance(serviceName string, endPoint string) {
	b := s.BalanceInstance[serviceName]
	b.Remove(endPoint)
}

func (s *Service) AddEndPointInBalance(serviceName string, endPoint string) {
	b := s.BalanceInstance[serviceName]
	b.Add(endPoint)
}

func (s *Service) Balance(serviceName string, key string) (string, error) {
	b := s.BalanceInstance[serviceName]
	return b.Balance(key)
}
