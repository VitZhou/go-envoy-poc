package cmd

import (
	"github.com/spf13/cobra"
	"go-envoy-poc/proxy"
	"go-envoy-poc/analyze"
	"net/http"
	"strconv"
	"log"
)

var configPath string
var RootCmd = &cobra.Command{
	Use:   "envoy",
	Short: "envoy",
	Long:  `envoy poc`,
	Run: func(cmd *cobra.Command, args []string) {
		newHttpProxy(configPath)
	},
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configPath, "configPath", "c", "./envoy.yaml", "envoy的配置文件路劲")
}

func newHttpProxy(path string) {
	staticResources, err := analyze.Parser(path)
	if err != nil {
		log.Fatalf("解析yaml文件错误%s", err)
	}

	h := proxy.NewHttpProxy(staticResources)
	port := staticResources.Address.Port
	if port <= 0 {
		log.Fatal("port必须大于0")
	}
	err = http.ListenAndServe(":"+strconv.Itoa(port), h)
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}
}
