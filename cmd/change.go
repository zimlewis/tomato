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
	Short: "Change the current tomato session",
	Long: `First argument is the session to change to(pomodoro, short, long)`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tomatoSession := args[0]
		var clockToSet int32
		switch strings.ToUpper(tomatoSession) {
		case "POMODORO": clockToSet = 0
		case "SHORT": clockToSet = 1
		case "LONG": clockToSet = 2
		default: 
			cmd.PrintErrf("Wrong clock: %s\nClock must be either pomodoro, short or long", tomatoSession)
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
