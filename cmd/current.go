package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/gen/proto/timer"
	errs "github.com/zimlewis/tomato/internal/tomatoerrs"
	"github.com/zimlewis/tomato/internal/formatter"
	"github.com/zimlewis/tomato/internal/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		// Print the time every second until the the program close
		for {
			<-ticker.C
			cur, err := getCurrentTime(ctx, c)
			if sta, ok := status.FromError(err); ok && (sta.Code() == codes.Canceled || sta.Code() == codes.Unavailable) {
				return
			}
			if errors.Is(err, errs.ErrCannotStopClock) {
				cmd.PrintErrln(err)
				break
			}
			if err != nil {
				cmd.PrintErrln(err)
				continue
			}
			s, err := f.Format(cur)
			cmd.Println(s)

			select {
			case <-ctx.Done(): return
			default:
			}
		}
	},
}

// Get the time and clock of the session
func getCurrentTime(ctx context.Context, c timer.TimerClient) (types.CurrentResponse, error) {
	// Get current time
	current, err := c.Current(ctx, nil)
	if sta, ok := status.FromError(err); ok && sta.Code() == codes.NotFound {
		// If the server return with not found, create a new session of current clock
		currentClock, err := c.GetClock(ctx, nil)

		if err != nil {
			return types.CurrentResponse{}, fmt.Errorf("Cannot get current clock: %w\n", err)
		}

		current = &timer.CurrentTimer{
			TimeLeft: int64(timeWait[currentClock.Clock] * 60),
			Clock: currentClock.Clock,
		}
		err = nil
	}
	if err != nil {
		return types.CurrentResponse{}, fmt.Errorf("Error while retrieving current time: %w\n", err)
	}

	remaining := current.TimeLeft

	// If the time is less than 0, stop the clock
	if remaining <= 0 {
		// If the clock cannot be stopped
		if _, err := c.Stop(ctx, nil); err != nil {
			return types.CurrentResponse{}, errors.Join(err, errs.ErrCannotStopClock)
		}

		// Send notification after successfully close the clock
		err = beeep.Notify("Tomato", "Your time is up", "")
		if err != nil {
			return types.CurrentResponse{}, fmt.Errorf("Cannot notify: %w", err)
		}

		// return the clock max time
		remaining = int64(timeWait[current.Clock]) * 60
	}


	return types.CurrentResponse{
		Clock: int16(current.Clock),
		TimeLeft: remaining,
	}, nil
}

func init() {
	rootCmd.AddCommand(currentCmd)

	currentCmd.Flags().StringP("formatter", "f", "default", "Decide which format to use")
}
