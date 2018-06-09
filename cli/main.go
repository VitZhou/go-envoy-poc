package main

import (
	"os"
	"go-envoy-poc/cli/cmd"
	"fmt"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
