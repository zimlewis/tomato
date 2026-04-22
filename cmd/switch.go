/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/storage"
)


// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		switch dir {
		case "up": switchUp() 
		case "down": switchDown()
		case "defualt": 
			fmt.Println("The argument to this command must be up or down")
		}
	},
}

func switchUp() {
	err := storage.Storage.Update(func(txn *badger.Txn) error {
		var i uint16

		item, err := txn.Get(timerKey)
		switch err {
		case badger.ErrKeyNotFound:
			i = 0
		case nil:
			b, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			i = binary.BigEndian.Uint16(b)
		default:
			return err
		}

		nextClock := i + 1
		if i >= 2 {
			nextClock = 0
		}


		nextClockByte := binary.BigEndian.AppendUint16(nil, nextClock)
		err = txn.Set(timerKey, nextClockByte)
		if err != nil {
			return err
		}

		err = txn.Delete(startTimeKey)

		return err
	})

	if err != nil {
		fmt.Printf("Something went wrong: %s\n", err.Error())
		os.Exit(1)
	}

}

func switchDown() {
	err := storage.Storage.Update(func(txn *badger.Txn) error {
		var i uint16

		item, err := txn.Get(timerKey)
		switch err {
		case badger.ErrKeyNotFound:
			i = 0
		case nil:
			b, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			i = binary.BigEndian.Uint16(b)
		default:
			return err
		}


		var nextClock uint16
		if i == 0 {
			nextClock = 2
		} else {
			nextClock = i - 1
		}


		nextClockByte := binary.BigEndian.AppendUint16(nil, nextClock)
		err = txn.Set(timerKey, nextClockByte)
		if err != nil {
			return err
		}

		err = txn.Delete(startTimeKey)

		return err
	})


	if err != nil {
		fmt.Printf("Something went wrong: %s\n", err.Error())
		os.Exit(1)
	}

}

func init() {
	rootCmd.AddCommand(switchCmd)
}
