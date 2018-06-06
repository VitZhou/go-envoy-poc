package cmd

import (
	"github.com/spf13/cobra"
	"go-envoy-poc/proxy"
	"go-envoy-poc/analyze"
	"net/http"
	"strconv"
	"go-envoy-poc/log"
)

var (
	configPath string
	logPath string
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

	h := proxy.NewHttpProxy(staticResources)
	port := staticResources.Address.Port
	if port <= 0 {
		log.Error.Fatal("port必须大于0")
	}
	err = http.ListenAndServe(":"+strconv.Itoa(port), h)
	if err != nil {
		log.Error.Fatalln("ListenAndServe:", err)
	}
}
