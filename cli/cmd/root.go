package cmd

import (
	"github.com/spf13/cobra"
	"go-envoy-poc/proxy"
	"go-envoy-poc/analyze"
	"net/http"
	"strconv"
	"go-envoy-poc/log"
	"go-envoy-poc/analyze/health_check/filter"
)

var (
	configPath string
	logPath    string
	RootCmd    = &cobra.Command{
		Use:   "envoy",
		Short: "envoy",
		Long:  `envoy poc`,
		Run: func(cmd *cobra.Command, args []string) {
			newHttpProxy(configPath)
		},
	}
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&configPath, "configPath", "c", "./envoy.yaml", "envoy的配置文件路劲")
	RootCmd.PersistentFlags().StringVarP(&logPath, "logPath", "l", "./envoy.log", "日志文件路径")
}

func newHttpProxy(path string) {
	staticResources, err := analyze.Parser(path)
	if err != nil {
		log.Error.Fatalf("解析yaml文件错误%s", err)
	}
	if len(staticResources.Clusters) <= 0 {
		log.Error.Fatal("至少配置一个集群")
	}

	port := staticResources.Address.Port
	if port <= 0 {
		log.Error.Fatal("port必须大于0")
	}
	if staticResources.Protocol == "dubbo" {
		proxy.NewSocketProxy(staticResources)
	} else {
		h := proxy.NewReverseProxy(staticResources)
		newHealthCheckFilter(staticResources)
		err = http.ListenAndServe(":"+strconv.Itoa(port), h)
		if err != nil {
			log.Error.Fatalln("ListenAndServe:", err)
		}
	}
}

func newHealthCheckFilter(staticResources *analyze.StaticResources) {
	check := staticResources.HealthCheck
	if check.Cluster != "" && check.Path != "" {
		for _, v := range staticResources.Clusters {
			if v.Name == check.Cluster {
				f := filter.Filter{Cluster: v, HealthCheck: check}
				f.Filter()
			}
		}
	}
}
