package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const Version string = "0.4.1"
const BuildDate string = "2025-10-01"

func cliCommandDisplayHelp(args []string) {
	displayVersion := StringInSlice("-v", args[1:]) || StringInSlice("--version", args[1:])

	if displayVersion {
		bold := color.New(color.Bold).SprintFunc()
		italic := color.New(color.Italic).SprintFunc()

		fmt.Println()
		fmt.Println(bold("⚡️ Git branch"))
		fmt.Println()
		fmt.Println("build date: ", bold(BuildDate))
		fmt.Println("version:        ", bold(Version))
		fmt.Println()
		fmt.Println("-v / --version    : display this help menu")
		fmt.Println("-d / --delete     : enters the delete mode")
		// fmt.Println("-c / --count <int>: how many branches to display. Cares about GIT_BRANCH_COUNT")
		fmt.Println()
		fmt.Println(italic("Need help?"))
		fmt.Println(italic("https://github.com/Tchoupinax/git-branch/issues"))
		os.Exit(0)
	}
}
