package load_balance

import (
	"go-envoy-poc/analyze/addr"
	"go-envoy-poc/log"
	"math/rand"
)

type RandomBalancer struct {
}

func (r *RandomBalancer) Balancing(addrs []addr.SocketAddress) *addr.Target {
	if addrs == nil || len(addrs) <= 0 {
		log.Error.Fatal("代理地址没有正确配置")
	}
	address := addrs[rand.Intn(len(addrs))]
	return &addr.Target{Host: address.Host, Port: address.Port}
}
