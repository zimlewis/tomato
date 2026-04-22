/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/storage"
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
		type currentBody struct {
			Clock string
			Minute int
			Second int
		}

		var current currentBody 
		err := storage.Storage.View(func(txn *badger.Txn) error {
			item, err := txn.Get(startTimeKey)
			switch err {
			case badger.ErrKeyNotFound:
				currentClockItem, err := txn.Get(timerKey)

				switch err {
				case nil:
				default: return err
				case badger.ErrKeyNotFound:
					current = currentBody {
						Clock: "Pomo",
						Minute: 25,
						Second: 00,
					}
					return nil
				}

				clockByte, err := currentClockItem.ValueCopy(nil)
				if err != nil {
					return err
				}

				clockInt := binary.BigEndian.Uint16(clockByte)

				current = currentBody {
					Clock: clock[clockInt],
					Minute: timeWait[clockInt],
					Second: 0,
				}


				return nil
			case nil:
			default:
				return err
			}

			b, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			i := binary.BigEndian.Uint64(b)
			
			currentClockItem, err := txn.Get(timerKey)
			if err != nil {
				return err
			}

			clockByte, err := currentClockItem.ValueCopy(nil)
			if err != nil {
				return err
			}

			clockInt := binary.BigEndian.Uint16(clockByte)
			elapsed := time.Now().Unix() - int64(i)


			total := timeWait[clockInt] * 60
			remaining := total - int(elapsed)
			

			minute := remaining / 60
			second := remaining % 60

			if minute <= 0 && second <= 0 {
				_ = beeep.Notify("Your time is up", "Move to your next phase", "")
				
				return storage.Storage.Update(func(txn *badger.Txn) error {
					return txn.Delete(startTimeKey)
				})
			}

			current = currentBody {
				Clock: clock[clockInt],
				Minute: int(minute),
				Second: int(second),
			}

			return nil
		})


		if err != nil {
			fmt.Printf("Cannot get current state %s\n", err.Error())
			os.Exit(1)
		}

		type returnBody struct {
			Text    string `json:"text"`
			Tooltip string `json:"tooltip"`
			Alt     string `json:"alt"`
			Class   string `json:"class"`
		}

		value := returnBody {
			Text: fmt.Sprintf("%.4s %02d:%02d", strings.ToUpper(current.Clock), current.Minute, current.Second),
			Class: "tomato",
			Tooltip: "+ Left click to start\n+ Right click to stop\n+ Scroll up to switch mod up\n+ Scroll down to switch mod down",
		}
		

		jsonData, err := json.Marshal(value)
		fmt.Println(string(jsonData))
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
