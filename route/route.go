package route

import "go-envoy-poc/analyze"

type Route interface {
	Filter(url string) *analyze.Cluster
}
