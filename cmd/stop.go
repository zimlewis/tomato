/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/storage"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := storage.Storage.Update(func(txn *badger.Txn) error {
			return txn.Delete(startTimeKey)
		})

		if err != nil {
			fmt.Printf("Cannot start: %s", err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

}
