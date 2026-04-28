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

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start tomato session with current phase",
	Long: `This work by saving the current time, and later, current method could subtract that saved
time and return the time left`,
	Args: cobra.NoArgs,



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
		_, err = c.Start(ctx, nil)
		if stas, ok := status.FromError(err); ok && stas.Code() == codes.Canceled {
			return
		}

		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("Start session successfully")
	},
}


func init() {
	rootCmd.AddCommand(startCmd)

	// startCmd.Flags().StringP("timer", "t", "short", "....")

	// startCmd.Flags().BoolP("short", "s", false, "Short break")
	// startCmd.Flags().BoolP("pomodoro", "p", false, "Pomodoro")
	// startCmd.Flags().BoolP("long", "l", false, "Long break")
}
