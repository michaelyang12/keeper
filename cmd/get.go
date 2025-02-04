/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/michaelyang12/keeper/db"
	"github.com/michaelyang12/keeper/logging"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an existing credential.",
	Long:  `Retrieve an existing credential by tag. The credential will be displayed. Use for login purposes. `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			logging.Error("Too many arguments provided. Usage: get <tag>")
			return
		}
		tag := args[0]
		cred, err := db.FetchExistingCredential(tag)
		if err != nil {
			logging.Error("Error getting credentials: %v\n", err)
		}
		logging.Success("Credentials retrieved for %s: \n", tag)
		logging.Highlight("Username: %s\nPassword: %s\n", cred.User, cred.Password)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
