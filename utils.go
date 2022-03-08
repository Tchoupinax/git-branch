package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"

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
	// fmt.Print("\033[H\033[2J") // Clear the terminal
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
