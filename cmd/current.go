/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/client"
	timer "github.com/zimlewis/tomato/gen/proto"
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

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		c := timer.NewTimerClient(conn.Connection)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			cur, err := printCurrentTimeInterval(ctx, c)
			if sta, ok := status.FromError(err); ok && sta.Code() == codes.Canceled {
				return
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			formatter := formatter.WaybarFormatter {}
			s, err := formatter.Format(cur)
			cmd.Println(s)

			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	},
}

func printCurrentTimeInterval(ctx context.Context, c timer.TimerClient) (types.CurrentResponse, error) {
	type returnBody struct {
		Text    string `json:"text"`
		Tooltip string `json:"tooltip"`
		Alt     string `json:"alt"`
		Class   string `json:"class"`
	}

	current, err := c.Current(ctx, nil)
	if sta, ok := status.FromError(err); ok && sta.Code() == codes.NotFound {
		currentClock, err := c.GetClock(ctx, nil)

		if err != nil {
			return types.CurrentResponse{}, fmt.Errorf("Cannot get current clock: %w\n", err)
		}

		return types.CurrentResponse{
			Clock: int16(currentClock.Clock),
			TimeLeft: int64(timeWait[current.Clock]),
		}, nil
	}
	if err != nil {
		return types.CurrentResponse{}, fmt.Errorf("Error while retrieving current time: %w\n", err)
	}

	remaining := current.TimeLeft

	if remaining <= 0 {

		err = beeep.Notify("Tomato", "Your time is up", "")
		if err != nil {
			return types.CurrentResponse{}, fmt.Errorf("Cannot notify: %w", err)
		}
		err = exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/complete.oga", "--volume=13076").Run()
		if err != nil {
			return types.CurrentResponse{}, fmt.Errorf("Cannot play sound: %w", err)
		}
		// Notify
		if _, err := c.Stop(ctx, nil); err != nil {
			return types.CurrentResponse{}, fmt.Errorf("Cannot stop the clock: %w\n", err)
		}

		remaining = int64(timeWait[current.Clock]) * 60
	}


	return types.CurrentResponse{
		Clock: int16(current.Clock),
		TimeLeft: remaining,
	}, nil
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
