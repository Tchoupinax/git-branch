package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	sf "github.com/sa-/slicefunk"
)

type SelectionBranch struct {
	Index             int
	Name              string
	Hash              string
	CommitedAt        time.Time
	CommitedAtTimeAgo string
}

func ChooseBranchNumber(branches []Branch) int {
	bold := color.New(color.Bold).SprintFunc()
	fmt.Println()
	fmt.Print(bold("✏️  Choose a branch : "))

	index := 0
	selectionBranches := sf.Map(branches, func(branch Branch) SelectionBranch {
		index++
		return SelectionBranch{
			Index:             index,
			Name:              branch.Name,
			Hash:              branch.Hash,
			CommitedAt:        branch.CommitedAt,
			CommitedAtTimeAgo: branch.CommitedAtTimeAgo,
		}
	})

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F31F {{ .Index | white }} {{ .Name | yellow }} ({{ .CommitedAtTimeAgo | red }})",
		Inactive: "   {{ .Index | white }} {{ .Name | cyan }} ({{ .CommitedAtTimeAgo | red }})",
		Selected: "\U0001F31F {{ .Index | white }} {{ .Name | red | cyan }}",
		Details: `
--------- Branches ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Commited At" | faint }}	{{ .CommitedAt }} ({{ .CommitedAtTimeAgo }})`,
	}

	searcher := func(input string, index int) bool {
		branch := selectionBranches[index]
		index = branch.Index
		name := strings.Replace(strings.ToLower(branch.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input) || input == strconv.Itoa(index)
	}

	prompt := promptui.Select{
		Label:     "Which branch to switch",
		Items:     selectionBranches,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Println(err)
	}

	return i
}
