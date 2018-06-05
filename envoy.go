package main

import (
	"net/http"
	"log"
	"go-envoy-poc/proxy"
)

func main() {
	h := &proxy.HttpProxy{DestHost: "127.0.0.1", DestPort: 9955}
	err := http.ListenAndServe(":9966", h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
