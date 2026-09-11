package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
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
		c, ctx, closeFunc, err := initializeClient()
		if err != nil {
			cmd.PrintErrf("error initializing client: %s", err.Error())
			return
		}
		defer closeFunc()

		// Get formatter using the "formatter" flag
		formatterFlag := cmd.Flag("formatter").Value.String()
		f := formatter.NewFromString(formatterFlag)

		stream, err := c.Current(ctx, nil)
		for {
			curr, err := stream.Recv();
			if errors.Is(err, io.EOF) { 
				cmd.Println("stream terminated")
				return 
			}
			if err != nil { break }
			
			if curr.TimeLeft == 0 {
				err := beeep.Notify("Tomato", "Your time is up", "")
				if err != nil {
					cmd.PrintErrf("error notifying: %v\n", err)
				}
				_, err = c.Stop(ctx, nil)
				if err != nil {
					cmd.PrintErrln(err)
					break
				}
				continue
			}


			current := &types.CurrentResponse{
				TimeLeft: int64(curr.TimeLeft),
				Clock: int16(curr.Clock),
			}
			s, err := f.Format(*current)
			fmt.Println(s)
		}
	},
}


func init() {
	rootCmd.AddCommand(currentCmd)

	currentCmd.Flags().StringP("formatter", "f", "default", "Decide which format to use")
}
