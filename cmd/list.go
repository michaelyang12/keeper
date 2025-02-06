/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/michaelyang12/keeper/db"
	"github.com/michaelyang12/keeper/logging"
	"github.com/michaelyang12/keeper/models"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all credentials.",
	Long: `Lists all existing credentials. Passwords are hidden by default for security reasons, but can be shown via the __ flag.
	Default format: <credential_tag> | <username>`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeList(args); err != nil {
			logging.Error("Failed to list credentials: %v\n", err)
		}
	},
}

func executeList(args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("invalid number of arguments")
	}

	credentials, err := db.FetchAllExistingCredentials()
	if err != nil {
		return fmt.Errorf("error listing credentials: %v", err)
	}
	printCredentialsList(credentials)
	return nil
}

func printCredentialsList(credentials []models.Credentials) {
	logging.Info("Stored credentials: \n")
	logging.Info("----------------------------------------------------------\n")
	logging.Highlight("%-15s | %-30s | %-20s\n", "Tag", "Username", "Password")
	logging.Info("----------------------------------------------------------\n")

	for _, cred := range credentials {
		logging.Display("%-15s | %-30s | %-20s\n", cred.Tag, cred.Username, cred.Password)
		// logging.PrintTabbedCredentials(cred.Tag, cred.Username, cred.Password)
	}

}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
