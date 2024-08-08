package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5/plumbing"
)

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func CheckIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GenerateSpace(length int) string {
	var str = "        "

	for i := 0; i < length; i++ {
		str += " "
	}

	return str
}

func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Criteria(r *plumbing.Reference, isRemote bool, isTag bool, displayAll bool) bool {
	if strings.Contains(string(r.Name()), "HEAD") {
		return false
	}

	if displayAll {
		return true
	}

	if isRemote {
		return r.Name().IsRemote()
	}

	if isTag {
		return r.Name().IsTag()
	}

	return !r.Name().IsRemote() && !r.Name().IsTag()
}

func getGitRootPath() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	stdout, err := cmd.Output()

	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(stdout))
}

func ChooseBranchNumber(placeholder string, deleteMode bool) int {
	input := read(placeholder, deleteMode)

	if !IsNumeric(input) {
		fmt.Println()
		red := color.New(color.Bold, color.BgHiRed).SprintFunc()
		fmt.Printf("%s", red("you should type a number"))

		fmt.Println()
		os.Exit(1)
	}

	intVar, _ := strconv.Atoi(input)
	return intVar - 1
}

func read(input string, deleteMode bool) string {
	bold := color.New(color.Bold).SprintFunc()
	red := color.New(color.Bold, color.BgHiRed).SprintFunc()

	fmt.Println()
	if deleteMode {
		fmt.Printf("%s%s", red(bold("✏️  Choose a branch : ")), input)
	} else {
		fmt.Printf("%s%s", bold("✏️  Choose a branch : "), input)
	}

	if err := keyboard.Open(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer func() {
		_ = keyboard.Close()
	}()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEnter {
			break
		}

		if char == 0 && key == 3 { // Ctrl + C
			keyboard.Close()
			fmt.Println("")

			os.Exit(0)
		}

		if char == 0 {
			if (len(input)) > 0 {
				input = input[:len(input)-1]
				if deleteMode {
					fmt.Printf("\033[2K\r%s%s", red(bold("✏️  Choose a branch : ")), input)
				} else {
					fmt.Printf("\033[2K\r%s%s", bold("✏️  Choose a branch : "), input)
				}
			}
		} else {
			fmt.Printf("%s", string(char))
			input = fmt.Sprintf("%s%s", input, string(char))
		}
	}

	return input
}
