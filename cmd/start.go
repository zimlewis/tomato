package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/gen/proto/timer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start tomato session with current phase",
	Long: `This work by saving the current time, and later, current method could subtract that saved
time and return the time left`,
	Args: cobra.NoArgs,

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

		_, err = c.Start(ctx, nil)
		if stas, ok := status.FromError(err); ok && stas.Code() == codes.Canceled { return }

		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("Start session successfully")
	},
}


func init() {
	rootCmd.AddCommand(startCmd)
}
