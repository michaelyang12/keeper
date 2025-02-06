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

var updateUsername bool

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update existing credentials",
	Long:  `Update credentials stored for the specified tag. By default, update takes one argument (password) to be updated. If you wish to update the username as well, the --user flag can be used.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeUpdate(args); err != nil {
			logging.Error("Failed to update credentials: %v\n", err)
		}
	},
	Example: "  keeper update google newPassword\n  keeper udpate google newPassword --user newemail@gmail.com",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("requires at least 2 arguments: <tag> <password>")
		} else if len(args) > 2 && !updateUsername {
			return fmt.Errorf("too many arguments; use --user to update username")
		}
		// } else if !updateUsername && len(args) < 3 {
		// 	return fmt.Errorf("missing argument; use --user to update username")
		// }
		return nil
	},
}

func executeUpdate(args []string) error {
	tag, password := args[0], args[1]

	var username string
	if updateUsername {
		username = args[2]
	} else {
		cred, err := db.FetchExistingCredential(tag)
		if err != nil {
			return err
		}
		username = cred.Username
	}
	if err := db.UpdateExistingCredential(tag, username, password); err != nil {
		return fmt.Errorf("error updating credentials: %w", err)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
