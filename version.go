package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const Version string = "0.0.11"
const BuildDate string = "2023-11-26"

func cliCommandDisplayVersion(args []string) {
	displayVersion := StringInSlice("-v", args[1:]) || StringInSlice("--version", args[1:])

	if displayVersion {
		bold := color.New(color.Bold).SprintFunc()
		fmt.Println()
		fmt.Println(bold("⚡️ Git branch"))
		fmt.Println()
		fmt.Println("build date: ", bold(BuildDate))
		fmt.Println("version:        ", bold(Version))
		fmt.Println()
		os.Exit(0)
	}
}
