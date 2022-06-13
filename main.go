package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"sync"

	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	// Check if the version is asked by flag
	cliCommandDisplayVersion(os.Args)

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
			branches = append(branches, Branch{Name: r.Name().Short(), Hash: r.Hash().String()})
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
		if len(branch.Name) > lengthOfGreatestBranchLength {
			lengthOfGreatestBranchLength = len(branch.Name)
		}
	}

	sort.Slice(branches, func(i, j int) bool {
		return branches[i].CommitedAt.After(branches[j].CommitedAt)
	})

	var isTheFirstArgumentIsANumber = (len(os.Args) > 1 && IsNumeric(os.Args[1])) || (len(os.Args) > 2 && (IsNumeric(os.Args[2])))

	bold := color.New(color.Bold).SprintFunc()
	red := color.New(color.Italic, color.FgRed).SprintFunc()
	if !isTheFirstArgumentIsANumber {
		fmt.Println("")
		fmt.Println(bold("⚡️ Git branch"))
	}

	var desiredBranchNumber int
	if isTheFirstArgumentIsANumber {
		tmpNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			tmpNumber, _ = strconv.Atoi(os.Args[2])
		}

		desiredBranchNumber = tmpNumber - 1
	} else {
		desiredBranchNumber = ChooseBranchNumber(branches)
	}

	if desiredBranchNumber > len(branches)-1 || desiredBranchNumber == -1 {
		fmt.Println(red("Please choose a number between 1 and ", len(branches)))

		os.Exit(1)
	}

	desiredBranch := branches[desiredBranchNumber]

	// worktree, _ := r.Worktree()
	// error := worktree.Checkout(&git.CheckoutOptions{
	// 	Branch: plumbing.NewBranchReferenceName(desiredBranch.name),
	// })

	// Temporary reclacement of the native command
	cmd := exec.Command("git", "checkout", desiredBranch.Name)
	cmd.Output()

	ClearTerminal()

	purple := color.New(color.Bold, color.FgHiMagenta).SprintFunc()
	fmt.Println(purple("Checkout the branch ", desiredBranch.Name, "!"))
}
