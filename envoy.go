package main

import (
	"net/http"
	"log"
	"net/http/httputil"
	"net/url"
	"math/rand"
)
func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		target := targets[rand.Int()%len(targets)]
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{Director: director}
}


func main(){
	h := &handle{host: "127.0.0.1", port: "9955"}
	err := http.ListenAndServe(":9966", h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
