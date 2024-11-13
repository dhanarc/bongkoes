package cmd

import (
	"fmt"
	"github.com/djk-lgtm/bongkoes/cmd/deployment"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deployment.InitCommand)
	rootCmd.AddCommand(deployment.GetConfigCommand)
	rootCmd.AddCommand(deployment.CreateCommand)
	rootCmd.AddCommand(deployment.GetLatestIssuesCommand)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
