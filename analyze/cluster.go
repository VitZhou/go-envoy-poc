package analyze

import (
	"go-envoy-poc/analyze/addr"
	"go-envoy-poc/load_balance"
	"sync"
	"go-envoy-poc/log"
)

type Cluster struct {
	Name   string
	Hosts  []addr.SocketAddress
	Policy string

	b            load_balance.Balancer
	validHosts   []addr.SocketAddress
	DeleteSignal chan addr.SocketAddress
	AddSignal    chan addr.SocketAddress
	m            sync.RWMutex
}

func (cluster *Cluster) Init() {
	cluster.AddSignal = make(chan addr.SocketAddress, 1)
	cluster.DeleteSignal = make(chan addr.SocketAddress, 1)
	cluster.initBalancer()
	cluster.m.RLock()
	cluster.validHosts = []addr.SocketAddress{}
	for _, v := range cluster.Hosts {
		cluster.validHosts = append(cluster.validHosts, addr.SocketAddress{Port: v.Port, Host: v.Host})
	}
	cluster.m.RUnlock()
	go func() {
		for {
			d := <-cluster.AddSignal
			cluster.addAddress(d)
		}
	}()

	go func() {
		for {
			d := <-cluster.DeleteSignal
			cluster.delAddress(d)
		}
	}()
}


func (cluster *Cluster) addAddress(add addr.SocketAddress) {
	if cluster.Exist(add) {
		return
	}
	cluster.m.Lock()
	cluster.validHosts = append(cluster.validHosts, add)
	log.Info.Println("add2", cluster.validHosts)
	cluster.m.Unlock()
}

func (cluster *Cluster) delAddress(del addr.SocketAddress) {
	if !cluster.Exist(del) {
		return
	}
	cluster.m.Lock()
	defer cluster.m.Unlock()
	for k, v := range cluster.validHosts {
		if v.Host == del.Host && v.Port == del.Port {
			cluster.validHosts = append(cluster.validHosts[:k], cluster.validHosts[k+1:]...)
			return
		}
	}
	log.Info.Println(cluster.validHosts)
}

func (cluster *Cluster) Exist(address addr.SocketAddress) bool {
	cluster.m.RLock()
	defer cluster.m.RUnlock()
	for _, v := range cluster.validHosts {
		if v.Host == address.Host && v.Port == address.Port {
			return true
		}
	}
	return false
}

func (cluster *Cluster) initBalancer() {
	switch cluster.Policy {
	case "round_robin":
		cluster.b = new(load_balance.RoundRobin)
		return
	default:
		cluster.b = new(load_balance.RoundRobin)
	}
}

func (cluster *Cluster) GetAddress() *addr.Target {
	cluster.m.RLock()
	log.Info.Println("当前存活实例:", cluster.validHosts)
	balancing := cluster.b.Balancing(cluster.validHosts)
	log.Info.Println("负载器筛选结果:", balancing)
	cluster.m.RUnlock()
	return balancing
}
