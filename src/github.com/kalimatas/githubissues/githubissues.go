// Generate HTML for Github issues

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"code.google.com/p/goauth2/oauth"
	"github.com/google/go-github/github"
)

const (
	issuesSeparator = ","
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
}
