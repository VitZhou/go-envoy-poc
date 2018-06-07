package analyze

import (
	"go-envoy-poc/load_balance"
	"go-envoy-poc/analyze/addr"
)

type Cluster struct {
	Name   string
	Hosts  []addr.SocketAddress
	Policy string
}

func (cluster *Cluster) GetAddress() *addr.Target {
	return load_balance.Balancing(cluster.Hosts)
}
