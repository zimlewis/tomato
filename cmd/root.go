/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)


var clock = []string{"pomodoro", "short", "long"}
var timerKey = []byte("timer")
var startTimeKey = []byte("start")
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
	Run: func(cmd *cobra.Command, args []string) { },
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

func init() {
}


