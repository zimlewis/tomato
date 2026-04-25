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

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Change the current tomato session",
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
		dir := args[0]

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		conn, err := client.New()
		if err != nil {
			cmd.PrintErrln(err)
		}
		defer func () {
			err := conn.Connection.Close()
			if err != nil {
				cmd.PrintErrln(err)
				return
			}
		}()

		c := timer.NewTimerClient(conn.Connection)

		
		switch dir {
		case "up": err = switchUp(ctx, c) 
		case "down": err = switchDown(ctx, c)
		case "defualt": 
			cmd.Println("The argument to this command must be up or down")
		}
		if stas, ok := status.FromError(err); ok && stas.Code() == codes.Canceled {
			return
		}
		if err != nil {
			cmd.PrintErrln(err)
		}
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
