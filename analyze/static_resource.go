package analyze

import (
	"go-envoy-poc/analyze/addr"
	"go-envoy-poc/analyze/health_check"
)

type StaticResources struct {
	Name    string
	Address addr.SocketAddress
	Routes   []Route
	Clusters []Cluster
	HealthCheck health_check.HttpHealthCheck `yaml:"health_check"`
}

type Route struct {
	Prefix string
	Cluster string
}



