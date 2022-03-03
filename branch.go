package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Branch struct {
	name       string
	hash       string
	commitedAt time.Time
}

func (b *Branch) GetCommitterDateFromLogs() {
	cmd := exec.Command("git", "show", b.name)
	stdout, _ := cmd.Output()

	var re = regexp.MustCompile(`(?m)Date:.*`)

	for index, match := range re.FindAllString(string(stdout), -1) {
		if index == 0 {
			var a = strings.Split(match, " +")[0]
			var c = strings.Split(a, "Date:   ")
			t, err := time.Parse(time.ANSIC, c[1])

			if err != nil {
				fmt.Println(err)
			}

			b.commitedAt = t.Add(-time.Hour * 1)
		}
	}
}
