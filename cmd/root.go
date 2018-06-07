package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func noArgs(cmd *cobra.Command, args []string) error {
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

func init() {
}

func Execute() {
	RootCmd.Execute()
}
