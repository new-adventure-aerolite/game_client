package cmd

import (
	"github.com/spf13/cobra"

	"github.com/new-adventure-aerolite/game-client/pkg/client"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a game",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client.Start()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
