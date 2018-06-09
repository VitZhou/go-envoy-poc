package load_balance

import (
	"sync/atomic"
	"go-envoy-poc/analyze/addr"
	"go-envoy-poc/log"
)

type Balancer interface {
	Balancing(addrs []addr.SocketAddress) *addr.Target
}


type RoundRobin struct {
	i int32
}

func (round *RoundRobin)Balancing(addrs []addr.SocketAddress) *addr.Target {
	if addrs == nil || len(addrs) <= 0{
		log.Error.Fatal("代理地址没有正确配置")
	}
	if atomic.LoadInt32(&round.i) >= int32(len(addrs)) {
		atomic.StoreInt32(&round.i, 0)
	}
	address := addrs[atomic.LoadInt32(&round.i)]
	atomic.AddInt32(&round.i, 1)
	return &addr.Target{Host: address.Host, Port: address.Port}
}
