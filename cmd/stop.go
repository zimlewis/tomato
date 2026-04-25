/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/client"
	timer "github.com/zimlewis/tomato/gen/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Reset current tomato session",
	Long: `Work by deleting saved time, if the current command cannot see the saved time, it will return maximum value for each session by default`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		conn, err := client.New()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		defer func () {
			err := conn.Connection.Close()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}()

		c := timer.NewTimerClient(conn.Connection)
		_, err = c.Stop(ctx, nil)
		if stas, ok := status.FromError(err); ok && stas.Code() == codes.Canceled {
			return
		}
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

}
