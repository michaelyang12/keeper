/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/michaelyang12/keeper/logging"
	"github.com/michaelyang12/keeper/utils"
	"github.com/spf13/cobra"
)

var length int

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a secure password",
	Long: `Generates a secure password using _ encryption. This can be used in conjunction with the add command to auto-generate a password for a new credentials.
	Usage:`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeGenerate(); err != nil {
			logging.Error("Failed to generate password: %v\n", err)
			return
		}
	},
	Args: cobra.ExactArgs(0),
}

func executeGenerate() error {
	password, err := utils.GenerateRandomPassphrase(length)
	if err != nil {
		return fmt.Errorf("error generating passphrase: %w", err)
	}

	err = clipboard.WriteAll(password)
	if err != nil {
		return fmt.Errorf("failed to copy to clipboard: %w", err)
	}

	logging.Success("Password generated and copied to clipboard!\n")
	logging.Highlight("%v\n", password)

	return nil
}

func init() {
	generateCmd.Flags().IntVarP(&length, "length", "l", 16, "Length of the generated password. Generated passwords must be 12 characters minimum for maximum security")
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
