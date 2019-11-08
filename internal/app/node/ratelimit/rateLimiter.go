/**
* @Author: HongBo Fu
* @Date: 2019/10/21 16:34
 */

package ratelimit

import "github.com/SongLiangChen/RateLimiter"

var rules RateLimiter.Rules
var rateLimiter RateLimiter.RateLimiter

type Limiter struct {
}

var LimiterInstance Limiter = Limiter{}

func init() {
	rules = RateLimiter.NewRules()
	//rateLimiter, _ = RateLimiter.NewRateLimiter("redis")
	rateLimiter, _ = RateLimiter.NewRateLimiter("memory")

}

func (l *Limiter) AddRule(limitKey string, rule *RateLimiter.Rule) {
	rules.AddRule(limitKey, rule)
}

func (l *Limiter) RemoveAllRule() {
	rules = map[string][]*RateLimiter.Rule{}
}

func (l *Limiter) Start() {
	//rateLimiter.InitRules(rules, "127.0.0.1:6379", "", "0", "10", "20")
	rateLimiter.InitRules(rules)
}

func (l *Limiter) GetLimiter() RateLimiter.RateLimiter {
	return rateLimiter
}
