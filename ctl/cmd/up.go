package cmd

import (
	"fmt"
	"github.com/k3rn3l-p4n1c/goaway/daemon/core"
	"github.com/spf13/cobra"
	"net/rpc"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {

		var (
			addr     = "/tmp/goaway.sock"
			request  = &core.Request{Name: ""}
			response = new(core.Response)
		)

		// Establish the connection to the adddress of the
		// RPC server
		client, _ := rpc.Dial("unix", addr)

		defer client.Close()

		// Perform a procedure call (core.HandlerName == Handler.Execute)
		// with the Request as specified and a pointer to a response
		// to have our response back.
		_ = client.Call(core.UpCmd, request, response)
		fmt.Println(response.Message)
	},
}

func init() {
	upCmd.PersistentFlags().Bool("bootstrap", false, "Run server as bootstrap")
	upCmd.PersistentFlags().String("join", "", "Try to join seed server")
}
