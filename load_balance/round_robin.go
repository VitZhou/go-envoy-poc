package load_balance

import (
	"sync/atomic"
	"sync"
	"go-envoy-poc/analyze/addr"
)

type Balancer interface {
	Balancing(addr []addr.SocketAddress) *addr.Target
}

var i int32 = 0
var m sync.Mutex

func Balancing(addrs []addr.SocketAddress) *addr.Target {
	m.Lock()
	if atomic.LoadInt32(&i) >= int32(len(addrs)) {
		atomic.StoreInt32(&i, 0)
	}
	m.Unlock()
	address := addrs[atomic.LoadInt32(&i)]
	atomic.AddInt32(&i, 1)
	return &addr.Target{Host: address.Host, Port: address.Port}
}
