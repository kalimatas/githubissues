// Generate HTML for Github issues

package main

import (
	"fmt"
	"os"
	"flag"
	"strings"
	"strconv"

	_"github.com/google/go-github/github"
)

const (
	issuesSeparator = ","
)

// Command-line arguments
var owner string
var repository string
var accessToken string
var milestone int
var issues string

func init() {
	flag.StringVar(&owner, "owner", "", "repository owner")
	flag.StringVar(&repository, "repository", "", "repository name")
	flag.StringVar(&accessToken, "access-token", "", "access token for authenticated access")
	flag.IntVar(&milestone, "milestone", 0, "milestone number")
	flag.StringVar(&issues, "issues", "", "comma-separated issue numbers list")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error: %v", r)
		}
	}()

	flag.Parse()

	if milestone == 0 && issues == "" {
		fmt.Fprint(os.Stderr, "Error: either milestone or issues must be set\n\n")
		flag.Usage()
		os.Exit(2)
	}

	var issueList []int
	if (issues != "") {
		for _, numberStr := range strings.Split(issues, issuesSeparator) {
			number, err := strconv.Atoi(numberStr)
			if err == nil {
				issueList = append(issueList, number)
			} else {
				panic(fmt.Sprintf("%v is not an issue number", numberStr))
			}
		}
	}
}
