package main

import (
    "os"
	"go-envoy-poc/cli/cmd"
)


func main(){
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
