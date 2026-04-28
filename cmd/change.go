/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"
	"os/signal"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/client"
	timer "github.com/zimlewis/tomato/gen/proto"
)

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tomatoSession := args[0]
		var clockToSet int32
		switch strings.ToUpper(tomatoSession) {
		case "POMODORO": clockToSet = 0
		case "SHORT": clockToSet = 1
		case "LONG": clockToSet = 2
		default: 
			cmd.PrintErrf("Wrong clock: %s\nClock must be eiter pomodoro, short or long", tomatoSession)
			return
		}

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
		
		_, err = c.SetClock(ctx, &timer.SetClockRequest{
			Clock: clockToSet,
		})

		cmd.Println("Change session successfully")
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)

}
