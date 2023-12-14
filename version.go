package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const Version string = "0.0.12"
const BuildDate string = "2023-12-14"

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
		fmt.Println(italic("Need help?"))
		fmt.Println(italic("https://github.com/Tchoupinax/git-branch/issues"))
		os.Exit(0)
	}
}
