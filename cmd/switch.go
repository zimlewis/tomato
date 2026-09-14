package cmd

import (
	"context"
	"errors"

	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/gen/proto/timer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Cycle the current tomato session",
	Long: `First argument is the direction which the command will change to
the commands is in this order: Pomodoro -> Short Break -> Long Break
eg.
from Pomodoro to Short Break:
	tomato switch up

from Long Break to Short Break:
	tomato switch down

Note that the operation will switch in a cycle, meaning it will go back to the first phase if it exceed last phase and vice versa
eg.
from Long Break to Pomodoro:
	tomato switch up

from Pomodoro to Long Break:
	tomato switch down
`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, ctx, closeFunc := initializeClient()
		defer closeFunc()

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
		
		dir := args[0]
		switch dir {
		case "up": err = switchUp(ctx, c) 
		case "down": err = switchDown(ctx, c)
		default: err = errors.New("The argument to this command must be up or down")
		}
		if stas, ok := status.FromError(err); ok && stas.Code() == codes.Canceled { return }
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("Switch session successfully")
	},
}

func switchUp(ctx context.Context, client timer.TimerClient) error {
	_, err := client.Switch(ctx, &timer.SwitchRequest{
		Dir: timer.SwitchDirection_UP,
	})
	return err
}

func switchDown(ctx context.Context, client timer.TimerClient) error {
	_, err := client.Switch(ctx, &timer.SwitchRequest{
		Dir: timer.SwitchDirection_DOWN,
	})
	return err
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
