package cmd

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

func initializeClient() (*client.Client, context.Context, func()) {
	port := os.Getenv("TOMATO_PORT")
	if port == "" {
		port = "6600"
	}
	host := "dns:///localhost"
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	timerClient := client.New(
		client.WithDialOptions(dialOptions...),
		client.WithHost(host),
		client.WithPort(port),
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(), 
		os.Interrupt,
	)


	return timerClient, ctx, cancel
}
