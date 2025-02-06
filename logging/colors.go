package logging

import (
	"github.com/fatih/color"
)

var (
	Error     = color.New(color.FgRed).PrintfFunc()
	Warn      = color.New(color.FgYellow).PrintfFunc()
	Success   = color.New(color.FgGreen).PrintfFunc()
	Info      = color.New(color.FgBlue).PrintfFunc()
	Display   = color.New(color.FgCyan).PrintfFunc()
	Highlight = color.New(color.FgHiMagenta).PrintfFunc()
)

func PrintTabbedCredentials(tag string, username string, password string) {
	// Highlight(fmt.Sprintf("%s\t", tag))
	// Display(fmt.Sprintf("%s\t", username))
	// Info(fmt.Sprintf("%s\t\n", password))
	Display("%s\t%s\t%s\n", tag, username, password)
	// fmt.Printf("%v\t%v\t%v\n", tag, username, password)
}
