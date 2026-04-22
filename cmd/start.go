package cmd

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/storage"
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
		err := storage.Storage.Update(func(txn *badger.Txn) error {
			_, err := txn.Get(timerKey)
			switch err {
			case badger.ErrKeyNotFound:
				err = txn.Set(timerKey, binary.BigEndian.AppendUint16(nil, 0))
			case nil:
			default:
				return err
			}

			currentTime := time.Now().Unix()
			bCurrent := binary.BigEndian.AppendUint64(nil, uint64(currentTime))

			err = txn.Set(startTimeKey, bCurrent)
			return err
		})

		if err != nil {
			fmt.Printf("Cannot start: %s", err.Error())
			os.Exit(1)
		}
	},
}


func init() {
	rootCmd.AddCommand(startCmd)

	// startCmd.Flags().StringP("timer", "t", "short", "....")

	// startCmd.Flags().BoolP("short", "s", false, "Short break")
	// startCmd.Flags().BoolP("pomodoro", "p", false, "Pomodoro")
	// startCmd.Flags().BoolP("long", "l", false, "Long break")
}
