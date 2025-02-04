package logging

import "github.com/fatih/color"

var (
	Error     = color.New(color.FgRed).PrintfFunc()
	Warn      = color.New(color.FgYellow).PrintfFunc()
	Success   = color.New(color.FgGreen).PrintfFunc()
	Info      = color.New(color.FgWhite).PrintfFunc()
	Highlight = color.New(color.FgCyan).PrintfFunc()
)
