package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,



	Run: func(cmd *cobra.Command, args []string) {
		t, err := cmd.Flags().GetString("timer")
		if err != nil {
			fmt.Println("Something went wrong, please try again later")
		}

		switch t {
		case "short": startShortBreak()
		case "long": startLongBreak()
		case "pomodoro": startPomodoro()
		default: return
		}
	},
}

func startPomodoro() {
	fmt.Println("Start pomodoro")
}

func startShortBreak() {
	fmt.Println("Start short break")
}

func startLongBreak() {
	fmt.Println("Start long break")
}


func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringP("timer", "t", "short", "....")

	// startCmd.Flags().BoolP("short", "s", false, "Short break")
	// startCmd.Flags().BoolP("pomodoro", "p", false, "Pomodoro")
	// startCmd.Flags().BoolP("long", "l", false, "Long break")
}
