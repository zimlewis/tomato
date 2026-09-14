package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/gen/proto/timer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Reset current tomato session",
	Long: `Work by deleting saved time, if the current command cannot see the saved time, it will return maximum value for each session by default`,
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx, cancel := initializeClient()
		defer cancel()

		connection, err := client.GetConnection()
		if err != nil {
			cmd.PrintErrf("cannot get client connection: %s", err.Error())
			return
		}
		defer func() {
			if err := connection.Close(); err != nil {
				cmd.PrintErrf("error during connection closing: %s", err.Error())
			}
		}()

		c := timer.NewTimerClient(connection)

		_, err = c.Stop(ctx, nil)
		if stas, ok := status.FromError(err); ok && stas.Code() == codes.Canceled { return }

		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("Stop timer successfully")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

}
