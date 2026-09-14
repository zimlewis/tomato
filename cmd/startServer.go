package cmd

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/dgraph-io/badger/v4"
	"github.com/spf13/cobra"
	"github.com/zimlewis/tomato/internal/badgerrepo"
	"github.com/zimlewis/tomato/internal/log"
	"github.com/zimlewis/tomato/internal/service/timer"
	"github.com/zimlewis/tomato/server"
)

// startServerCmd represents the startServer command
var startServerCmd = &cobra.Command{
	Use:   "ss",
	Short: "Start the gRPC server that whole connection to the local Database",
	Long: `Since it is not possible to have multiple connection to the badger db at the same time
the server will hold the database connection and then will retrieve send data to client(other cli command)
by gRPC protocal`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		// Initialize dependency of the server
		logger, closeFunc, err := log.InitializeLogger()
		if err != nil {
			cmd.PrintErrf("cannot initialize logger: %s", err.Error())
		}
		defer func() {
			err := closeFunc()
			if err != nil { cmd.PrintErrf("error closing logger: %s", err.Error()) }
		}()

		// Get the environment variable
		tomatoStorage := os.Getenv("TOMATO_STORAGE")
		if tomatoStorage == "" { tomatoStorage = "/tmp/tomato/" }
		tomatoPort := os.Getenv("TOMATO_PORT")
		if tomatoPort == "" { tomatoPort = "6600" }

		// Open badger database
		badgerOptions := badger.DefaultOptions(tomatoStorage)
		badgerDB, err := badger.Open(badgerOptions)
		if err != nil {
			logger.Error(
				"cannot open badger database",
				slog.Any("error", err),
			)
		}
		defer badgerDB.Close()

		repo := badgerrepo.New(
			badgerrepo.WithDatabase(badgerDB),
		)
		service := timer.New(
			timer.WithLogger(logger),
			timer.WithRepository(&repo),
		)

		// Create new server with dependency created and start the server
		timerServer := server.New(
			server.WithHost("0.0.0.0"),
			server.WithPort(tomatoPort),
			server.WithLogger(logger),
			server.WithService(service),
		)

		err = timerServer.Start(ctx)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("closing server")
	},
}

func init() {
	rootCmd.AddCommand(startServerCmd)
}
