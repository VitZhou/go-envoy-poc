package route

import (
	"go-envoy-poc/analyze"
	"strings"
	"go-envoy-poc/log"
)

type PrefixRoute struct {
	Routes     []analyze.Route
	clusterMap map[string]*analyze.Cluster
}

func NewPrefixRoute(routes []analyze.Route, clusters []analyze.Cluster) *PrefixRoute {
	route := &PrefixRoute{Routes: routes}
	route.clusterMap = make(map[string]*analyze.Cluster)
	for k := range clusters {
		route.clusterMap[clusters[k].Name] = &clusters[k]
	}
	return route
}

func (prefixRoute *PrefixRoute) Filter(url string) *analyze.Cluster {
	for _, v := range prefixRoute.Routes {
		if strings.HasPrefix(url, v.Prefix) {
			cluster, exists := prefixRoute.clusterMap[v.Cluster]
			if !exists {
				log.Error.Printf("路由规则没有相匹配的集群,集群%s", v.Cluster)
				return nil
			}
			return cluster
		}
	}
	return nil
}
