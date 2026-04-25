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
	Use:   "ss",
	Short: "Start the gRPC server that whole connection to the local Database",
	Long: `Since it is not possible to have multiple connection to the badger db at the same time
the server will hold the database connection and then will retrieve send data to client(other cli command)
by gRPC protocal`,
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
