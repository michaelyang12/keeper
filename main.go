/*
Copyright © 2025 Michael Yang <michaelyang17@outlook.com>
*/
package main

import (
	"github.com/michaelyang12/keeper/cmd"
	"github.com/michaelyang12/keeper/db"
	"github.com/michaelyang12/keeper/logging"
	"github.com/michaelyang12/keeper/utils"
)

// var rootCmd = &cobra.Command{
// 	Use:   "keeper",
// 	Short: "A simple CLI password manager. ",
// }

func main() {
	if err := db.InitializeLocalDatabase(); err != nil {
		logging.Error("Error at root: %v\n", err)
		return
	}

	// Store key (wont overrite if already exists)
	if err := utils.StoreKey(); err != nil {
		logging.Error("Error storing encryption key :%v\n", err)
		return
	}

	cmd.Execute()
}
