package route

import (
	"testing"
	"go-envoy-poc/analyze"
)

func TestFilter(t *testing.T) {
	t.Run("正确匹配", func(t *testing.T) {
		routes := []analyze.Route{{Prefix: "/", Cluster: "clusterName"}}
		clusters := []analyze.Cluster{{Name: "clusterName", Host: "localhost", Port: 80}}
		route := NewPrefixRoute(routes, clusters)

		target := route.Filter("/1")

		if target.Port != 80 || target.Host != "localhost" {
			t.Error("fail")
		}
	})

	t.Run("没有匹配的路由", func(t *testing.T) {
		routes := []analyze.Route{{Prefix: "/", Cluster: "clusterName"}}
		clusters := []analyze.Cluster{{Name: "clusterName", Host: "localhost", Port: 80}}
		route := NewPrefixRoute(routes, clusters)

		target := route.Filter("aa/")

		if target != nil {
			t.Error("fail")
		}
	})

	t.Run("匹配了路径,但是找不到集群", func(t *testing.T) {
		routes := []analyze.Route{{Prefix: "/", Cluster: "clusterName2"}}
		clusters := []analyze.Cluster{{Name: "clusterName", Host: "localhost", Port: 80}}
		route := NewPrefixRoute(routes, clusters)

		target := route.Filter("/")

		if target != nil {
			t.Error("fail")
		}
	})
}
