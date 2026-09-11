package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/client"
	"github.com/zimlewis/tomato/gen/proto/timer"
)

var timeWait = []int{25, 5, 30}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tomato",
	Short: "A CLI tool to track time with tomato timer",
	Long: `A CLI tool to track time with tomato timer

tomato is a CLI tool that track time with tomato timer. It has 3 mod:
 + Pomodoro: 25 minutes of working
 + Short break: 5 minutes of break
 + Long break: 30 minutes of break

it is recommended to do this rotation while working/studying:
(Pomodoro -> Short break) x 5 -> Long break
and then repeat
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetErr(os.Stderr)
	rootCmd.SetOut(os.Stdout)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func initializeClient() (timer.TimerClient, context.Context, func() error, error) {
	// Create new connection to the server
	conn, err := client.New()
	if err != nil {
		return nil, nil, nil, err
	}

	// Create new context that caught os.Interrupt event so closing the client would be handled gracefully
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	c := timer.NewTimerClient(conn.Connection)

	// Create close function that close the connection and cancel the context
	closeFunc := func() error {
		cancel()
		return conn.Connection.Close()
	}

	return c, ctx, closeFunc, nil
}
