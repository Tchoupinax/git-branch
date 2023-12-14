package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func cliCommandDisplayVersion(args []string) {
	displayVersion := StringInSlice("-h", args[1:]) || StringInSlice("--help", args[1:])

	if displayVersion {
		bold := color.New(color.Bold).SprintFunc()

		fmt.Println()
		fmt.Println(bold("⚡️ Git branch"))
		fmt.Println()
		fmt.Println("-h/--help:          displays this menu")
		fmt.Println()
		os.Exit(0)
	}
}
