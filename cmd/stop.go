package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Reset current tomato session",
	Long: `Work by deleting saved time, if the current command cannot see the saved time, it will return maximum value for each session by default`,
	Run: func(cmd *cobra.Command, args []string) {
		c, ctx, closeFunc, err := initializeClient()
		if err != nil {
			cmd.PrintErrf("error initializing client: %s", err.Error())
			return
		}
		defer closeFunc()

		_, err = c.Stop(ctx, nil)
		if stas, ok := status.FromError(err); ok && stas.Code() == codes.Canceled { return }

		if err != nil {
			cmd.PrintErrln(err)
			return
		}

		cmd.Println("Stop timer successfully")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

}
