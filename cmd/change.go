package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/gen/proto/timer"
)

// changeCmd represents the change command
var changeCmd = &cobra.Command{
	Use:   "change",
	Short: "Change the current tomato session",
	Long: `First argument is the session to change to (pomodoro, short, long)`,

	Args: cobra.ExactArgs(1),
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


		// Get the correct clock base on the first argument of the command
		tomatoSession := args[0]
		var clockToSet int32
		switch strings.ToUpper(tomatoSession) {
		case "POMODORO": clockToSet = 0
		case "SHORT": clockToSet = 1
		case "LONG": clockToSet = 2
		default: 
			// Print an error if the argument isn't correct
			cmd.PrintErrf("Wrong clock: %s\nClock must be either pomodoro, short or long", tomatoSession)
			return
		}

		_, err = c.SetClock(ctx, &timer.SetClockRequest{
			Clock: clockToSet,
		})
		if err != nil {
			cmd.PrintErrf("Error setting clock: %s", err.Error())
			return
		}

		// No error, successfully change the clock
		cmd.Println("Change session successfully")
	},
}

func init() {
	rootCmd.AddCommand(changeCmd)

}
