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
	commitedAt time.Time
	hash       string
}

func (b *Branch) GetCommitterDateFromLogs() {
	cmd := exec.Command("git", "show", b.name)
	stdout, _ := cmd.Output()

	var re = regexp.MustCompile(`(?m)Date:.*`)

	for index, match := range re.FindAllString(string(stdout), -1) {
		if index == 0 {
			a := strings.Split(match, " +")[0]

			if strings.Contains(a, " -") {
				a = strings.Split(match, " -")[0]
			}

			var c = strings.Split(a, "Date:   ")

			t, err := time.Parse(time.ANSIC, c[1])

			if err != nil {
				fmt.Println(err)
			}

			// *1 in winter and *2 in summer... >.<
			b.commitedAt = t.Add(-time.Hour * 2)
		}
	}
}
