package analyze

import "go-envoy-poc/analyze/addr"

type StaticResources struct {
	Name    string
	Address addr.SocketAddress
	Routes   []Route
	Clusters []Cluster
}

type Route struct {
	Prefix string
	Cluster string
}



