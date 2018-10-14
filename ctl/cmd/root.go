package cmd

import (
	"fmt"

	"github.com/k3rn3l-p4n1c/goaway/daemon/core"
	"github.com/k3rn3l-p4n1c/goaway/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/rpc"
)

func noArgs(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return fmt.Errorf(
		"goaway: '%s' is not a goaway command.\nSee 'goaway --help'", args[0])
}

var RootCmd = &cobra.Command{
	Use:           "goaway [OPTIONS] COMMAND [ARG...]",
	Short:         "Go away.",
	Long:          `GO AWAY! engine. An orchestration scheduler for microservices.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Args:          noArgs,
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Print value",
	Run: func(cmd *cobra.Command, args []string) {

		var (
			addr     = "/tmp/goaway.sock"
			request  = &core.Request{Command: ""}
			response = new(core.Response)
		)
		// Establish the connection to the adddress of the
		// RPC server
		client, err := rpc.Dial("unix", addr)

		defer client.Close()
		if err != nil {
			logrus.WithError(err).Fatal("unable to connect to goaway daemon")
		}

		// Perform a procedure call (core.HandlerName == Handler.Execute)
		// with the Request as specified and a pointer to a response
		// to have our response back.
		err = client.Call(core.GetCmd, request, response)
		if err != nil {
			logrus.WithError(err).Fatal("unable to call to goaway daemon")
		}
		fmt.Println("Get:", response.Message)

	},
}

var setCmd = &cobra.Command{
	Use:   "set [NAME]",
	Short: "Set value",
	Run: func(cmd *cobra.Command, args []string) {

		var (
			addr     = "/tmp/goaway.sock"
			request  = &core.Request{Command: args[0]}
			response = new(core.Response)
		)
		// Establish the connection to the adddress of the
		// RPC server
		client, err := rpc.Dial("unix", addr)

		defer client.Close()
		if err != nil {
			logrus.WithError(err).Fatal("unable to connect to goaway daemon")
		}

		// Perform a procedure call (core.HandlerName == Handler.Execute)
		// with the Request as specified and a pointer to a response
		// to have our response back.
		err = client.Call(core.SetCmd, request, response)
		if err != nil {
			logrus.WithError(err).Fatal("unable to call to goaway daemon")
		}
		fmt.Println("Get:", response.Message)

	},
}

func init() {
	RootCmd.AddCommand(upCmd)
	RootCmd.AddCommand(getCmd)
	RootCmd.AddCommand(setCmd)
	RootCmd.AddCommand(version.Cmd)
	RootCmd.AddCommand(stackCmd)

	viper.BindPFlag("master.bootstrap", upCmd.PersistentFlags().Lookup("bootstrap"))
	viper.BindPFlag("slave.join", upCmd.Flags().Lookup("join"))
}

func Execute() {
	RootCmd.Execute()
}
