/**
* @Author: HongBo Fu
* @Date: 2019/10/16 17:08
 */

package cfg

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"sync"
)

var (
	Instance ConfigInstance
)

type ConfigInstance struct {
	Config Config
	sync.Mutex
}

type Config struct {
	Routes       []Route `json:"routes" yaml:"routes"`
	Port         string  `json:"port" yaml:"port"`
	ReadTimeout  int64   `json:"readTimeout" yaml:"readTimeout"`
	WriteTimeout int64   `json:"writeTimeout" yaml:"writeTimeout"`
}

type Route struct {
	Route                  string   `json:"route" yaml:"route"`
	NLBModel               string   `json:"loadBalanceModel" yaml:"loadBalanceModel"`
	LoadBalanceKey         string   `json:"loadBalanceKey" yaml:"loadBalanceKey"`
	ProcessPlugin          Process  `json:"processPlugin" yaml:"processPlugin"`
	HealthCheck            bool     `json:"healthCheck" yaml:"healthCheck"`
	Swagger                string   `json:"swagger" yaml:"swagger"`
	CheckHealthURL         string   `json:"checkHealthURL" yaml:"checkHealthURL"`
	CheckTicket            int      `json:"checkTicket" yaml:"checkTicket"`
	Endpoints              []string `json:"endpoints" yaml:"endpoints"`
	HealthLimitTimeWindow  int32    `json:"healthLimitTimeWindow" yaml:"healthLimitTimeWindow"`
	HealthLimitCount       int32    `json:"healthLimitCount" yaml:"healthLimitCount"`
	RequestLimitTimeWindow int32    `json:"requestLimitTimeWindow" yaml:"requestLimitTimeWindow"`
	RequestLimitCount      int32    `json:"RequestLimitCount" yaml:"RequestLimitCount"`
}

type Process struct {
	InBound    []Plugin `json:"inBound" yaml:"inBound"`
	OutBound   []Plugin `json:"outBound" yaml:"outBound"`
	ErrorBound []Plugin `json:"errorBound" yaml:"errorBound"`
}

type Plugin struct {
	ActionType string                 `json:"actionType" yaml:"actionType"`
	Properties map[string]interface{} `json:"properties" yaml:"properties"`
}

func parseJsonFile(configFile string) Config {
	var config = Config{}
	data, err := ioutil.ReadFile(configFile)
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func parseYamlFile(configFile string) Config {
	var config = Config{}
	data, err := ioutil.ReadFile(configFile)
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func Reload(configFile string, fileType string) {
	Instance.Lock()
	defer Instance.Unlock()
	var t Config
	if fileType == "json" {
		t = parseJsonFile(configFile)
	} else {
		t = parseYamlFile(configFile)
	}

	oldConfig := Instance.Config

	if !reflect.DeepEqual(t, oldConfig) {
		Instance.Config = t
	}
}
