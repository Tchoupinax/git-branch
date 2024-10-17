package main

import (
	"fmt"
	math "math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"

	timeago "github.com/caarlos0/timea.go"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	// Check if the version is asked by flag
	cliCommandDisplayVersion(os.Args)
	// Check if the helper is asked by flag
	cliCommandDisplayHelp(os.Args)

	var deleteMode = StringInSlice("-d", os.Args[1:]) || StringInSlice("--delete", os.Args[1:])
	var branchCount float64 = 10
	if os.Getenv("GIT_BRANCH_COUNT") != "" {
		value, err := strconv.Atoi(os.Getenv("GIT_BRANCH_COUNT"))
		if err != nil {
			CheckIfError(err)
		} else {
			branchCount = float64(value)
		}
	}

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

	err := refIter.ForEach(func(r *plumbing.Reference) error {
		var displayAll = StringInSlice("-a", os.Args[1:]) || StringInSlice("--all", os.Args[1:])
		var displayRemoteBranches = StringInSlice("-r", os.Args[1:])
		var displayTags = StringInSlice("-t", os.Args[1:])

		if Criteria(r, displayRemoteBranches, displayTags, displayAll) {
			branches = append(branches, Branch{name: r.Name().Short(), hash: r.Hash().String()})
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

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

	sort.Slice(branches, func(i, j int) bool {
		return branches[i].commitedAt.After(branches[j].commitedAt)
	})

	var isTheFirstArgumentIsANumber = (len(os.Args) > 1 && IsNumeric(os.Args[1])) || (len(os.Args) > 2 && (IsNumeric(os.Args[2])))

	bold := color.New(color.Bold).SprintFunc()
	yellow := color.New(color.Bold, color.FgYellow).SprintFunc()
	blue := color.New(color.Bold, color.FgBlue).SprintFunc()
	red := color.New(color.Italic, color.FgRed).SprintFunc()

	if !isTheFirstArgumentIsANumber {
		fmt.Println("")
		if deleteMode {
			fmt.Println(bold("⚡️ Git branch"), red(bold("DELETE MODE")))
		} else {
			fmt.Println(bold("⚡️ Git branch"))
		}
		fmt.Println("")

		var count = 1
		for branchIndex, branch := range branches[:int(math.Min(float64(len(branches)), branchCount))] {
			var s string
			if branch.commitedAt.String() != "0001-01-01 00:00:00 +0000 UTC" {
				s = timeago.Of(branch.commitedAt)
			} else {
				s = "—"
			}

			space := " "
			if branchIndex < 9 {
				space = "  "
			}

			if count%2 == 0 {
				fmt.Println(blue(count, space, branch.name, GenerateSpace(lengthOfGreatestBranchLength-len(branch.name)), "       (", s, ")"))
			} else {
				fmt.Println(yellow(count, space, branch.name, GenerateSpace(lengthOfGreatestBranchLength-len(branch.name)), "       (", s, ")"))
			}

			count = count + 1
		}
	}

	var desiredBranchNumber int
	if isTheFirstArgumentIsANumber {
		tmpNumber, err := strconv.Atoi(os.Args[1])
		if err != nil {
			tmpNumber, _ = strconv.Atoi(os.Args[2])
		}

		desiredBranchNumber = tmpNumber - 1
	} else {
		desiredBranchNumber = ChooseBranchNumber("", deleteMode)
	}

	if desiredBranchNumber > len(branches)-1 || desiredBranchNumber == -1 {
		fmt.Println(red("Please choose a number between 1 and ", len(branches)))

		os.Exit(1)
	}

	// When option -a/--all is used to display remote branches,
	// we want to remove origin/ prefix to have local branches on the computer
	desiredBranchName := strings.Replace(branches[desiredBranchNumber].name, "origin/", "", 1)

	if deleteMode {
		// Temporary reclacement of the native command
		cmd := exec.Command("git", "branch", "-D", desiredBranchName)
		var _, cmdError = cmd.Output()
		if cmdError != nil {
			panic(cmdError)
		}

		ClearTerminal()

		cyan := color.New(color.Bold, color.FgCyan).SprintFunc()
		bgYellow := color.New(color.Bold, color.BgYellow).SprintFunc()
		fmt.Println(bgYellow(cyan("branch ", desiredBranchName, " deleted!")))
	} else {
		// worktree, _ := r.Worktree()
		// error := worktree.Checkout(&git.CheckoutOptions{
		// 	Branch: plumbing.NewBranchReferenceName(desiredBranch.name),
		// })

		// Temporary reclacement of the native command
		cmd := exec.Command("git", "checkout", desiredBranchName)
		var _, cmdError = cmd.Output()
		if cmdError != nil {
			panic(cmdError)
		}

		ClearTerminal()

		purple := color.New(color.Bold, color.FgHiMagenta).SprintFunc()
		fmt.Println(purple("Checkout the branch ", desiredBranchName, "!"))
	}
}
