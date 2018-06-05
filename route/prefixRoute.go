package route

import (
	"go-envoy-poc/analyze"
	"strings"
	"log"
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

func (this *PrefixRoute) Filter(url string) *Target {
	for _, v := range this.Routes {
		if strings.HasPrefix(url, v.Prefix) {
			cluster,exists := this.clusterMap[v.Cluster]
			if !exists{
				log.Fatalf("路由规则没有相匹配的集群,集群%s", v.Cluster)
				return nil
			}
			return &Target{Host: cluster.Host, Port: cluster.Port}
		}
	}
	return nil
}
