/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/michaelyang12/keeper/db"
	"github.com/michaelyang12/keeper/logging"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new credential",
	Long: `Add a new credential to Keeper. keeper
	Usage: add <tag> <username> <pasword>`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			fmt.Println("Error: Invalid number of arguments \nCorrect usage: add <tag> <username> <password>.")
			return
		}

		// Get credentials
		tag := args[0]
		username := args[1]
		password := args[2]

		// TODO: Store credentials
		if err := db.InsertNewCredential(tag, username, password); err != nil {
			logging.Error("Error storing credentials: %v\n", err)
			return
		}

		// Print new credentials
		logging.Success("Added new credential: \nTag: %s\nUsername: %s\nPassword: %s\n", tag, username, password)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
