package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start tomato session with current phase",
	Long: `This work by saving the current time, and later, current method could subtract that saved
time and return the time left`,
	Args: cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		c, ctx, closeFunc, err := initializeClient()
		if err != nil {
			cmd.PrintErrf("error initializing client: %s", err.Error())
			return
		}
		defer closeFunc()

		_, err = c.Start(ctx, nil)
		if stas, ok := status.FromError(err); ok && stas.Code() == codes.Canceled { return }

		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("Start session successfully")
	},
}


func init() {
	rootCmd.AddCommand(startCmd)
}
