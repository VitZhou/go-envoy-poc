package proxy

import (
	"net/http"
	"net/url"
	"net/http/httputil"
	"strconv"
	"log"
	"strings"
	"go-envoy-poc/analyze"
	"go-envoy-poc/route"
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
		log.Fatal("路由配置错误")
	}
	remote, err := url.Parse("http://" + target.Host + ":" + strconv.Itoa(target.Port))
	if err != nil {
		log.Fatalf("创建代理失败%s", err)
	}
	proxy := newSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}

func newSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {

		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
