package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func CheckIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func clearTerminal() {
	fmt.Print("\033[H\033[2J") // Clear the terminal
}

func criteria(r *plumbing.Reference, isRemote bool, isTag bool) bool {
	if strings.Contains(string(r.Name()), "HEAD") {
		return false
	}

	if isRemote && isTag {
		return r.Name().IsRemote() && r.Name().IsTag()
	}

	if isRemote {
		return r.Name().IsRemote() || !r.Name().IsRemote()
	}

	if isTag {
		return r.Name().IsTag()
	}

	return !r.Name().IsRemote()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	path, _ := os.Getwd()

	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{})
	if err != nil {
		red := color.New(color.Bold, color.BgHiRed).SprintFunc()
		fmt.Println(red(path, " ", "is not a git repository :/"))

		os.Exit(1)
	}

	clearTerminal()
	branches := []string{}
	refIter, _ := r.References()
	refIter.ForEach(func(r *plumbing.Reference) error {
		if criteria(r, stringInSlice("-a", os.Args[1:]), false) {
			branches = append(branches, r.Name().Short())
		}

		return nil
	})

	bold := color.New(color.Bold).SprintFunc()
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()
	blue := color.New(color.Bold, color.FgBlue).SprintFunc()
	fmt.Println(bold("⚡️ Git branch"))
	fmt.Println("")

	var count = 1
	for _, v := range branches {
		if count%2 == 0 {
			fmt.Println(blue(count, " ", v))
		} else {
			fmt.Println(yellow(count, " ", v))
		}

		count = count + 1
	}

	fmt.Println()
	fmt.Print(bold("✏️  Choose a branch : "))

	input := read()
	fmt.Println()

	if !isNumeric(input) {
		fmt.Println()
		red := color.New(color.Bold, color.BgHiRed).SprintFunc()
		fmt.Printf(red("you should type a number"))

		fmt.Println()
		os.Exit(1)
	}

	worktree, _ := r.Worktree()
	intVar, _ := strconv.Atoi(input)
	desiredBranch := branches[intVar-1]

	fmt.Println(desiredBranch, plumbing.NewBranchReferenceName(desiredBranch))

	error := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(desiredBranch),
	})

	fmt.Println(error)

	purple := color.New(color.Bold, color.FgHiMagenta).SprintFunc()
	// clearTerminal()
	fmt.Println(purple("Checkout the branch ", desiredBranch, "!"))
}

func read() string {
	var input string

	if err := keyboard.Open(); err != nil {
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

			os.Exit(0)
		}

		fmt.Printf("%s", string(char))
		input = fmt.Sprintf("%s%s", input, string(char))

	}

	return input
}
