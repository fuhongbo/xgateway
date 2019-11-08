/**
* @Author: HongBo Fu
* @Date: 2019/10/17 10:40
 */

package node

import (
	"errors"
	"github.com/lafikl/liblb/bounded"
	"github.com/lafikl/liblb/consistent"
	"github.com/lafikl/liblb/p2c"
	"github.com/lafikl/liblb/r2"
)

type Balance struct {
	algorithm     interface{}
	algorithmName string
}

func (b *Balance) New(algorithmName string) {
	b.algorithmName = algorithmName
	switch algorithmName {
	case "RoundRobin":
		b.algorithm = r2.New()
	case "Random":
		b.algorithm = p2c.New()
	case "ConsistentHashing":
		b.algorithm = consistent.New()
	case "BoundedConsistentHashing":
		b.algorithm = bounded.New()
	}
}

func (b *Balance) Add(server string) {
	switch b.algorithmName {
	case "RoundRobin":
		al := b.algorithm.(*r2.R2)
		al.Add(server)
		b.algorithm = al
	case "Random":
		al := b.algorithm.(*p2c.P2C)
		al.Add(server)
		b.algorithm = al
	case "ConsistentHashing":
		al := b.algorithm.(*consistent.Consistent)
		al.Add(server)
		b.algorithm = al
	case "BoundedConsistentHashing":
		al := b.algorithm.(*bounded.Bounded)
		al.Add(server)
		b.algorithm = al
	}
}

func (b *Balance) Remove(server string) {
	switch b.algorithmName {
	case "RoundRobin":
		al := b.algorithm.(*r2.R2)
		al.Remove(server)
		b.algorithm = al
	case "Random":
		al := b.algorithm.(*p2c.P2C)
		al.Remove(server)
		b.algorithm = al
	case "ConsistentHashing":
		al := b.algorithm.(*consistent.Consistent)
		al.Remove(server)
		b.algorithm = al
	case "BoundedConsistentHashing":
		al := b.algorithm.(*bounded.Bounded)
		al.Remove(server)
		b.algorithm = al
	}
}

func (b *Balance) Balance(key string) (string, error) {
	switch b.algorithmName {
	case "RoundRobin":
		al := b.algorithm.(*r2.R2)
		return al.Balance()
	case "Random":
		al := b.algorithm.(*p2c.P2C)
		return al.Balance(key)
	case "ConsistentHashing":
		al := b.algorithm.(*consistent.Consistent)
		return al.Balance(key)
	case "BoundedConsistentHashing":
		al := b.algorithm.(*bounded.Bounded)
		return al.Balance(key)
	}
	return "", errors.New("没有找到对应的负载均衡算法")
}
