package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	timeago "github.com/caarlos0/timea.go"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	rootpath := getGitRootPath()

	if rootpath == "" {
		red := color.New(color.Bold, color.BgHiRed).SprintFunc()
		currentPath, _ := os.Getwd()
		fmt.Println(red(currentPath, " ", "is not a git repository :/"))

		os.Exit(1)
	}

	r, _ := git.PlainOpenWithOptions(rootpath, &git.PlainOpenOptions{})

	ClearTerminal()
	branches := []Branch{}
	refIter, _ := r.References()

	refIter.ForEach(func(r *plumbing.Reference) error {
		var displayAll = StringInSlice("-a", os.Args[1:]) || StringInSlice("--all", os.Args[1:])
		var displayRemoteBranches = StringInSlice("-r", os.Args[1:])
		var displayTags = StringInSlice("-t", os.Args[1:])

		if Criteria(r, displayRemoteBranches, displayTags, displayAll) {
			branches = append(branches, Branch{name: r.Name().Short(), hash: r.Hash().String()})
		}

		return nil
	})

	var wg sync.WaitGroup
	for i := range branches {
		wg.Add(1)

		go func(b *Branch) {
			defer wg.Done()
			b.GetCommitterDateFromLogs()
		}(&branches[i])
	}
	wg.Wait()

	ClearTerminal()

	var lengthOfGreatestBranchLength = 0
	for _, branch := range branches {
		if len(branch.name) > lengthOfGreatestBranchLength {
			lengthOfGreatestBranchLength = len(branch.name)
		}
	}

	var isTheFirstArgumentIsANumber = (len(os.Args) > 1 && IsNumeric(os.Args[1])) || (len(os.Args) > 2 && (IsNumeric(os.Args[2])))

	bold := color.New(color.Bold).SprintFunc()
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()
	blue := color.New(color.Bold, color.FgBlue).SprintFunc()
	red := color.New(color.Italic, color.FgRed).SprintFunc()
	if !isTheFirstArgumentIsANumber {
		fmt.Println("")
		fmt.Println(bold("⚡️ Git branch"))
		fmt.Println("")

		var count = 1
		for _, branch := range branches {
			s := timeago.Of(branch.commitedAt)

			if count%2 == 0 {
				fmt.Println(blue(count, "  ", branch.name, GenerateSpace(lengthOfGreatestBranchLength-len(branch.name)), "       (", s, ")"))
			} else {
				fmt.Println(yellow(count, "  ", branch.name, GenerateSpace(lengthOfGreatestBranchLength-len(branch.name)), "       (", s, ")"))
			}

			count = count + 1
		}
	}

	worktree, _ := r.Worktree()

	var desiredBranchNumber int
	if isTheFirstArgumentIsANumber {
		tmpNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			tmpNumber, _ = strconv.Atoi(os.Args[2])
		}

		desiredBranchNumber = tmpNumber - 1
	} else {
		desiredBranchNumber = chooseBranchNumber()
	}

	desiredBranch := branches[desiredBranchNumber]
	error := worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(desiredBranch.name),
	})

	ClearTerminal()

	if error != nil {
		fmt.Println(red(error))
		fmt.Println("")
	}

	purple := color.New(color.Bold, color.FgHiMagenta).SprintFunc()
	fmt.Println(purple("Checkout the branch ", desiredBranch.name, "!"))
}

func chooseBranchNumber() int {
	bold := color.New(color.Bold).SprintFunc()
	fmt.Println()
	fmt.Print(bold("✏️  Choose a branch : "))

	input := read()
	fmt.Println()

	if !IsNumeric(input) {
		fmt.Println()
		red := color.New(color.Bold, color.BgHiRed).SprintFunc()
		fmt.Printf(red("you should type a number"))

		fmt.Println()
		os.Exit(1)
	}

	intVar, _ := strconv.Atoi(input)
	return intVar - 1
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
			fmt.Println("")

			os.Exit(0)
		}

		fmt.Printf("%s", string(char))
		input = fmt.Sprintf("%s%s", input, string(char))

	}

	return input
}
