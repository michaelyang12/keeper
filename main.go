/*
Copyright © 2025 Michael Yang <michaelyang17@outlook.com>
*/
package main

import (
	"github.com/michaelyang12/keeper/cmd"
	"github.com/michaelyang12/keeper/db"
	"github.com/michaelyang12/keeper/logging"
)

// var rootCmd = &cobra.Command{
// 	Use:   "keeper",
// 	Short: "A simple CLI password manager. ",
// }

func main() {
	if err := db.InitializeLocalDatabase(""); err != nil {
		logging.Error("Error at root: %v\n", err)
		return
	}
	cmd.Execute()
}
