package proxy

import (
	"net/http"
	"net/url"
	"strconv"
	"go-envoy-poc/analyze"
	"go-envoy-poc/route"
	"go-envoy-poc/log"
	"net/http/httputil"
)

type HttpProxy struct {
	StaticResources *analyze.StaticResources
	route           route.Route
}


func NewHttpProxy(resources *analyze.StaticResources) *HttpProxy {
	routes := resources.Routes
	clusters := resources.Clusters
	prefixRoute := route.NewPrefixRoute(routes, clusters)
	return &HttpProxy{StaticResources: resources, route: prefixRoute}
}

func (httpProxy *HttpProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	target := httpProxy.route.Filter(path)
	if target == nil{
		log.Error.Fatal("路由配置错误")
	}
	remote, err := url.Parse("http://" + target.Host + ":" + strconv.Itoa(target.Port))
	if err != nil {
		log.Error.Fatalf("创建代理失败%s", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

