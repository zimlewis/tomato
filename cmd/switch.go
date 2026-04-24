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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		conn, err := client.New()
		if err != nil {
			cmd.PrintErrln(err)
		}

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
