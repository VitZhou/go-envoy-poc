package proxy

import (
	"net/http"
	"net/url"
	"net/http/httputil"
	"strconv"
	"log"
)

type HttpProxy struct {
	TargetHost string
	TargetPort int
}

func (httpProxy *HttpProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	remote, err := url.Parse("http://" + httpProxy.TargetHost + ":" + strconv.Itoa(httpProxy.TargetPort))
	if err != nil {
		log.Fatalf("创建代理失败%s", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.ServeHTTP(w, r)
}
