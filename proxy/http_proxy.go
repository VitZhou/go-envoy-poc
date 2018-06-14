package proxy

import (
	"net/http"
	"net/url"
	"strconv"
	"go-envoy-poc/analyze"
	"go-envoy-poc/route"
	"go-envoy-poc/log"
	"net/http/httputil"
	"strings"
	"time"
	"net"
)

type HttpProxy struct {
	StaticResources *analyze.StaticResources
	route           route.Route
}

func NewReverseProxy(resources *analyze.StaticResources) *httputil.ReverseProxy  {
	routes := resources.Routes
	clusters := resources.Clusters
	for k := range clusters {
		clusters[k].Init()
	}
	prefixRoute := route.NewPrefixRoute(routes, clusters)
	reverseProxy := HttpProxy{StaticResources: resources, route: prefixRoute}
	return reverseProxy.newSingleHostReverseProxy()
}

func (httpProxy *HttpProxy)newSingleHostReverseProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		path := req.URL.Path
		cluster := httpProxy.route.Filter(path)
		if cluster == nil {
			log.Error.Fatal("路由配置错误")
		}
		address := cluster.GetAddress()
		target, err := url.Parse("http://" + address.Host + ":" + strconv.Itoa(address.Port))
		if err != nil {
			log.Error.Fatalf("创建代理失败%s", err)
		}
		targetQuery := target.RawQuery

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
	return &httputil.ReverseProxy{Director: director, Transport: EnvoyTransport}
}

var EnvoyTransport http.RoundTripper = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
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
