// Generate HTML for Github issues

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

const (
	issuesSeparator = ","
	halfString = "&frac12;"

	mainTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset='utf-8'>
	<style>
		* {
			font-family: Arial, sans-serif;
		}
		.issue {
			text-align: left;
			width: 700px;
			border: 1px solid black;
			border-collapse: collapse;
			margin-bottom: 20px;
		}
		.left {
			width: 70px;
			text-align: center;
			vertical-align: top;
		}
		.big {
			font-size: 25px;
		}
		.description {
			min-height: 150px;
		}
		th, td {
			padding: 6px;
			border: 1px solid black;
		}
	</style>
</head>
<body>

{{range .}}

<table class="issue">
	<thead>
		<tr>
			<th class="left big">{{.Number}}</th>
			<th class="big">{{.Title}}</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td></td>
			<td><div class="description">{{.Body}}</div></td>
		</tr>
		<tr>
			<td class="left big">2SP</td>
			<td></td>
		</tr>
	</tbody>
</table>

{{end}}

</body>
</html>
`
)

// Command-line arguments
var owner string
var repository string
var accessToken string
var milestone string
var issues string

type IssueManager struct {
	client     *github.Client
	owner      string
	repository string
}

// fetchIssues fetches issues by milestone from Github tracker
func (m *IssueManager) fetchByMilestone(milestone string) ([]github.Issue, error) {
	issues, _, err := m.client.Issues.ListByRepo(m.owner, m.repository, &github.IssueListByRepoOptions{Milestone: milestone})
	if err != nil {
		return nil, err
	}

	return issues, err
}

// fetchIssues fetches issues by numbers from Github tracker
func (m *IssueManager) fetchByNumbers(issueList []int) ([]github.Issue, error) {
	var issues []github.Issue

	for _, issueNumber := range issueList {
		issue, _, err := m.client.Issues.Get(m.owner, m.repository, int(issueNumber))
		if err != nil {
			return issues, err
		}

		issues = append(issues, *issue)
	}

	return issues, nil
}

func (m *IssueManager) printHtml(issues []github.Issue) (err error) {
	template, err := template.New("issues").Parse(mainTemplate)
	if err != nil {
		return
	}

	template.Execute(os.Stdout, issues)

	return
}

// newClient returns *github.Client
//
// With empty accessToken a client with not authorized access is created (for public repos).
func newClient(accessToken string) (client *github.Client) {
	if accessToken == "" {
		client = github.NewClient(nil)
	} else {
		transport := &oauth.Transport{
			Token: &oauth.Token{AccessToken: accessToken},
		}
		client = github.NewClient(transport.Client())
	}

	return
}

func init() {
	flag.StringVar(&owner, "owner", "", "repository owner")
	flag.StringVar(&repository, "repository", "", "repository name")
	flag.StringVar(&accessToken, "access-token", "", "access token for authenticated access")
	flag.StringVar(&milestone, "milestone", "", "milestone number")
	flag.StringVar(&issues, "issues", "", "comma-separated issue numbers list")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error: %v", r)
		}
	}()

	flag.Parse()

	if milestone == "" && issues == "" {
		fmt.Fprint(os.Stderr, "Error: either milestone or issues must be set\n\n")
		flag.Usage()
		os.Exit(2)
	}

	var issueList []int
	if issues != "" {
		for _, numberStr := range strings.Split(issues, issuesSeparator) {
			number, err := strconv.Atoi(numberStr)
			if err == nil {
				issueList = append(issueList, number)
			} else {
				panic(fmt.Sprintf("%v is not an issue number", numberStr))
			}
		}
	}

	manager := &IssueManager{
		client: newClient(accessToken),
		owner: owner,
		repository: repository,
	}

	var issues []github.Issue
	var fetchError error
	if len(issueList) > 0 {
		issues, fetchError = manager.fetchByNumbers(issueList)
	} else {
		issues, fetchError = manager.fetchByMilestone(milestone)
	}

	if fetchError != nil {
		panic(fmt.Sprintf("cannot fetch issues: %v", fetchError))
	}

	err := manager.printHtml(issues)
	if err != nil {
		panic(fmt.Sprintf("cannot print html: %v", err))
	}
}
