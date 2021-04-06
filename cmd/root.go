package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "game-client",
	Short: "A Game Client that used to interact with Game Server via REST API",
	Long: `
	████████╗██╗░░██╗███████╗  ██╗░░██╗██╗███╗░░██╗░██████╗░  ░█████╗░███████╗  ██████╗░██████╗░
	╚══██╔══╝██║░░██║██╔════╝  ██║░██╔╝██║████╗░██║██╔════╝░  ██╔══██╗██╔════╝  ██╔══██╗██╔══██╗
	░░░██║░░░███████║█████╗░░  █████═╝░██║██╔██╗██║██║░░██╗░  ██║░░██║█████╗░░  ██║░░██║██████╦╝
	░░░██║░░░██╔══██║██╔══╝░░  ██╔═██╗░██║██║╚████║██║░░╚██╗  ██║░░██║██╔══╝░░  ██║░░██║██╔══██╗
	░░░██║░░░██║░░██║███████╗  ██║░╚██╗██║██║░╚███║╚██████╔╝  ╚█████╔╝██║░░░░░  ██████╔╝██████╦╝
	░░░╚═╝░░░╚═╝░░╚═╝╚══════╝  ╚═╝░░╚═╝╚═╝╚═╝░░╚══╝░╚═════╝░  ░╚════╝░╚═╝░░░░░  ╚═════╝░╚═════╝░`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
