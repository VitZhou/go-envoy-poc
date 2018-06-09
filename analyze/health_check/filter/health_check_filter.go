package filter

import (
	"go-envoy-poc/analyze"
	"time"
	"io/ioutil"
	"net/http"
	"go-envoy-poc/analyze/addr"
	"strconv"
	"go-envoy-poc/analyze/health_check"
)

type Filter struct {
	Cluster     analyze.Cluster
	HealthCheck health_check.HttpHealthCheck
}

func (f *Filter) Filter() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			for _, v := range f.Cluster.Hosts {
				go func(address addr.SocketAddress, check health_check.HttpHealthCheck) {
					resp, err := http.Get("http://" + address.Host + ":" + strconv.Itoa(address.Port) + check.Path)
					if err != nil {
						f.Cluster.DeleteSignal <- address
						return
					}
					defer resp.Body.Close()
					_, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						f.Cluster.DeleteSignal <- address
						return
					}
					f.Cluster.AddSignal <- address
				}(v, f.HealthCheck)
			}
		}
	}()
}
