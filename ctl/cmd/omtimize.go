package cmd

import (
	"github.com/k3rn3l-p4n1c/goaway/scheduler"
	"github.com/spf13/cobra"
)

var optimizeCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Run optimization",
	Run: func(cmd *cobra.Command, args []string) {
		scheduler.Run()
	},
}
