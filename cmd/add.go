/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/michaelyang12/keeper/db"
	"github.com/michaelyang12/keeper/logging"
	"github.com/michaelyang12/keeper/utils"
	"github.com/spf13/cobra"
)

var generate bool

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:     "add [flags] <tag> <username> <password>",
	Short:   "Add a new credential",
	Long:    `Add a new credential to Keeper.`,
	Example: `  keeper add google john@gmail.com password123`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeAdd(args); err != nil {
			logging.Error("Failed to add credentials: %v\n", err)
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("requires at least 2 arguments: <tag> <user> [password]")
		}
		if generate && len(args) > 2 {
			return fmt.Errorf("cannot specify both a password argument and --generate")
		}
		if !generate && len(args) < 3 {
			return fmt.Errorf("missing password argument; use --generate or provide a password")
		}
		return nil
	},
}

func executeAdd(args []string) error {
	// Get credentials
	tag := args[0]
	username := args[1]
	var password string
	if generate {
		p, err := utils.GenerateRandomPassphrase(16)
		if err != nil {
			return fmt.Errorf("error generating random password: %w", err)
		}
		password = p
		logging.Success("Random password generated!\n")
	} else {
		password = args[2]
	}

	// TODO: Store credentials
	if err := db.InsertNewCredential(tag, username, password); err != nil {
		return fmt.Errorf("error storing credentials: %v", err)
	}

	// Print new credentials
	logging.Success("Added new credentials: \n")
	logging.Display("Tag: %s\nUsername: %s\nPassword: %s\n", tag, username, password)
	return nil
}

func init() {
	addCmd.Flags().BoolVarP(&generate, "generate", "g", false, "Generate a random password instead of providing one")
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
