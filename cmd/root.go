/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "keeper",
	Short: "A simple CLI password manager made with Golang",
	Long: `Keeper is a simple CLI password manager made with Golang. It uses AES encryption and stores everything on the local machine securely using sqlcipher.
	To add credentials you wish to keep track of, run 'keeper add <tag> <username> <password>'.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// fmt.Println("Initializing db...")
	// 	// Set password

	// 	// fmt.Println("Successfully initialized local db")
	// },
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.keeper.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
