package cmd

import (
	"github.com/spf13/cobra"
	"github.com/k3rn3l-p4n1c/goaway/daemon/core"
	"net/rpc"
	"fmt"
)

var stackCmd = &cobra.Command{
	Use:   "stack COMMAND",
	Short: "Stack",
}
var stackUp = &cobra.Command{
	Use:   "up",
	Short: "Create stack",
	Run: func(cmd *cobra.Command, args []string) {

		var (
			addr     = "/tmp/goaway.sock"
			request  = &core.Request{Command: "./file/dummy.yml"}
			response = new(core.Response)
		)

		// Establish the connection to the adddress of the
		// RPC server
		client, _ := rpc.Dial("unix", addr)

		defer client.Close()

		// Perform a procedure call (core.HandlerName == Handler.Execute)
		// with the Request as specified and a pointer to a response
		// to have our response back.
		_ = client.Call(core.StackUpCmd, request, response)
		if response.Error {
			fmt.Printf("unable to start stack, err=%s\n", response.Message)
		}
		fmt.Println("Stack Up!!!")
		fmt.Println(response.Message)
	},
}

var stackOptimize = &cobra.Command{
	Use:   "optimize",
	Short: "Optimize stack",
	Run: func(cmd *cobra.Command, args []string) {

		var (
			addr     = "/tmp/goaway.sock"
			request  = &core.Request{Command: "./file/dummy.yml"}
			response = new(core.Response)
		)

		// Establish the connection to the adddress of the
		// RPC server
		client, _ := rpc.Dial("unix", addr)

		defer client.Close()

		// Perform a procedure call (core.HandlerName == Handler.Execute)
		// with the Request as specified and a pointer to a response
		// to have our response back.
		_ = client.Call(core.StackOptimizeCmd, request, response)
		if response.Error {
			fmt.Printf("unable to optimize stack, err=%s\n", response.Message)
		}
		fmt.Println("Stack Optimize!!!")
		fmt.Println(response.Message)
	},
}

func init() {
	stackCmd.AddCommand(stackUp)
	stackCmd.AddCommand(stackOptimize)
}
