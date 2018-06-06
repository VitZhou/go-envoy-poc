package route

import (
	"go-envoy-poc/analyze"
	"strings"
	"go-envoy-poc/log"
)

type PrefixRoute struct {
	Routes     []analyze.Route
	Clusters   []analyze.Cluster
	clusterMap map[string]analyze.Cluster
}

func NewPrefixRoute(routes []analyze.Route, clusters []analyze.Cluster) *PrefixRoute {
	route := &PrefixRoute{Routes: routes, Clusters: clusters}
	route.clusterMap = make(map[string]analyze.Cluster)
	for _, v := range clusters {
		route.clusterMap[v.Name] = v
	}
	return route
}

func (prefixRoute *PrefixRoute) Filter(url string) *Target {
	for _, v := range prefixRoute.Routes {
		if strings.HasPrefix(url, v.Prefix) {
			cluster,exists := prefixRoute.clusterMap[v.Cluster]
			if !exists{
				log.Error.Printf("路由规则没有相匹配的集群,集群%s", v.Cluster)
				return nil
			}
			return &Target{Host: cluster.Host, Port: cluster.Port}
		}
	}
	return nil
}
