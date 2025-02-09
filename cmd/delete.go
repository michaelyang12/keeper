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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing credential",
	Long: `Delete an existing credential.
	Usages: > --delete <username>
			> --delete <password>`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeDelete(args); err != nil {
			logging.Error("Failed to delete credentials: %v\n", err)
		}
	},
}

func executeDelete(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of arguments")
	}
	tag := args[0]

	if err := db.DeleteExistingCredentials(tag); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
