/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/server"
	"github.com/zimlewis/tomato/storage"
)

// startServerCmd represents the startServer command
var startServerCmd = &cobra.Command{
	Use:   "startServer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		storage.Initialize()
		err := server.Start(ctx)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(startServerCmd)
}
