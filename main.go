package main

import (
	"fmt"
	"os"
	"github.com/k3rn3l-p4n1c/goaway/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
