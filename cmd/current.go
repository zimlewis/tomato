package cmd

import (
	"errors"
	"io"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/gen/proto/timer"
	"github.com/zimlewis/tomato/internal/formatter"
	"github.com/zimlewis/tomato/internal/types"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Print the current time you have left in the phase and phase",
	Long: `Print each second current time in the phase, in waybar module json format
example output: 
{"text":"POMO 25:00","tooltip":"+ Left click to start\n+ Right click to stop\n+ Scroll up to switch mod up\n+ Scroll down to switch mod down","alt":"","class":"tomato"}
	`,

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

		// Get formatter using the "formatter" flag
		formatterFlag := cmd.Flag("formatter").Value.String()
		f := formatter.NewFromString(formatterFlag)

		stream, err := c.Current(ctx, nil)
		if err != nil {
			cmd.PrintErrf("cannot get the stream from server: %v", err)
			return
		}
		for {
			// Receive data from the stream
			curr, err := stream.Recv();
			// If server closed, close the session
			if errors.Is(err, io.EOF) { 
				cmd.Println("stream terminated")
				return 
			}
			if err != nil {
				cmd.PrintErrf("error retrieving from server: %v\n", err)
				break 
			}
			
			// If time left is 0, notify and then stop the clock and continue with next loop
			if curr.TimeLeft == 0 {
				err := beeep.Notify("Tomato", "Your time is up", "")
				// If notifying failed, print the error and ignore it since it's optional anyway
				if err != nil {
					cmd.PrintErrf("error notifying: %v\n", err)
				}
				// If stopping failed, break the loop because it will spam notification otherwise
				_, err = c.Stop(ctx, nil)
				if err != nil {
					cmd.PrintErrln(err)
					break
				}
				continue
			}

			// Format the current response and print it
			current := &types.CurrentResponse{
				TimeLeft: int64(curr.TimeLeft),
				Clock: int16(curr.Clock),
			}
			s, err := f.Format(*current)
			cmd.Println(s)
		}
	},
}


func init() {
	rootCmd.AddCommand(currentCmd)

	currentCmd.Flags().StringP("formatter", "f", "default", "Decide which format to use")
}
