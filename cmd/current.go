/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/client"
	timer "github.com/zimlewis/tomato/gen/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		conn, err := client.New()
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		c := timer.NewTimerClient(conn.Connection)
		for {
			s, err := printCurrentTimeInterval(ctx, c)
			if sta, ok := status.FromError(err); ok && sta.Code() == codes.Canceled {
				return
			}
			if err != nil {
				cmd.PrintErrln(err)
				return
			}

			cmd.Println(s)
			time.Sleep(1 * time.Second)
		}
	},
}

func printCurrentTimeInterval(ctx context.Context, c timer.TimerClient) (string, error) {
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
			return "", fmt.Errorf("Cannot get current clock: %w\n", err)
		}

		value := returnBody{
			Text: fmt.Sprintf(
				"%.4s %02d:%02d",
				strings.ToUpper(clock[currentClock.Clock]),
				timeWait[currentClock.Clock],
				0,
			),
			Class:   "tomato",
			Tooltip: "+ Left click to start\n+ Right click to stop\n+ Scroll up to switch mod up\n+ Scroll down to switch mod down",
		}

		jsonData, err := json.Marshal(value)
		if err != nil {
			return "", fmt.Errorf("Cannot marshal data: %w", err)
		}
		return fmt.Sprintf("%s", string(jsonData)), nil
	}
	if err != nil {
		return "", fmt.Errorf("Error while retrieving current time: %w\n", err)
	}

	remaining := current.TimeLeft

	if remaining <= 0 {
		// Notify
		if _, err := c.Stop(ctx, nil); err != nil {
			return "", fmt.Errorf("Cannot stop the clock: %w\n", err)
		}

		remaining = 0
	}

	minute := remaining / 60
	second := remaining % 60

	value := returnBody{
		Text: fmt.Sprintf(
			"%.4s %02d:%02d",
			strings.ToUpper(clock[current.Clock]),
			minute,
			second,
		),
		Class:   "tomato",
		Tooltip: "+ Left click to start\n+ Right click to stop\n+ Scroll up to switch mod up\n+ Scroll down to switch mod down",
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("Cannot marshal data: %w", err)
	}
	return fmt.Sprintf("%s", string(jsonData)), nil
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
