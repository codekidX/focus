package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/go-github/v32/github"
)

// GHRepo stores the repository info
type GHRepo struct {
	ID          int
	Name        string
	Description string
	FullName    string
	Private     bool
	HTMLURL     string
	OpenIssues  int
}

// GHIssue is the github issues response object
type GHIssue struct {
	URL           string `json:"url,omitempty"`
	HTMLURL       string `json:"html_url,omitempty"`
	RepositoryURL string `json:"repository_url,omitempty"`
	ID            int    `json:"id,omitempty"`
	Number        int    `json:"number,omitempty"`
	Title         string `json:"title,omitempty"`
	State         string `json:"state,omitempty"`
	// Milestone     GHMilestone `json:"milestone,omitempty"`
	// Assignee GHUser
	ClosedAt string `json:"closed_at,omitempty"`
	Body     string `json:"body,omitempty"`
}

// ListIssues returns the list of GithubIssues
func ListIssues(queryStr string) ([]GHIssue, error) {
	var issues []GHIssue
	base, err := GetAPIURL()
	if err != nil {
		return issues, err
	}
	url := base + "/issues"
	if queryStr != "" {
		url += queryStr
	}

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return issues, errors.New("OXO:Error while getting github issues")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return issues, err
	}

	err = json.Unmarshal(b, &issues)
	if err != nil {
		return issues, err
	}
	return issues, nil
}

// GetIssue gets a issue by id from github
func GetIssue(number int) (GHIssue, error) {
	var issue GHIssue
	base, err := GetAPIURL()
	if err != nil {
		return issue, err
	}

	url := base + fmt.Sprintf("/issues/%d", number)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return issue, errors.New("OXO:Error while getting github issue")
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return issue, err
	}

	err = json.Unmarshal(b, &issue)
	if err != nil {
		return issue, err
	}
	return issue, nil

}

// GetRepositoryURL gets the current directory repository url
func GetRepositoryURL() (string, error) {
	cmd := exec.Command("git", "remote", "-v")
	b, err := cmd.Output()
	if err != nil {
		return "", err
	}
	split := strings.Split(string(b), "\n")
	splitReplaced := strings.ReplaceAll(split[0], " ", ":")
	lineSplit := strings.Split(splitReplaced, ":")
	if len(lineSplit) != 3 {
		return "", errors.New("something went wrong while getting remote")
	}
	return lineSplit[1], nil
}

// GetAPIURL returns the api url for accessing github data
func GetAPIURL() (string, error) {
	owner, repo, err := GetRepOwnerAndName()
	if err != nil {
		return "", err
	}
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	if strings.HasSuffix(apiURL, ".git") {
		dotIndex := strings.LastIndex(apiURL, ".")
		apiURL = apiURL[:dotIndex]
	}
	return apiURL, nil
}

// GetRepOwnerAndName gets the owner and repo name of the
// current git directory
func GetRepOwnerAndName() (string, string, error) {
	repoURL, err := GetRepositoryURL()
	if err != nil {
		return "", "", err
	}

	split := strings.Split(repoURL, "/")
	owner := split[len(split)-2]
	repo := split[len(split)-1]
	if strings.HasSuffix(repo, ".git") {
		dotIndex := strings.LastIndex(repo, ".")
		repo = repo[:dotIndex]
	}
	return owner, repo, nil
}

// GetRepositoryInfo returns the current repository info from the
// github api
func GetRepositoryInfo() (GHRepo, error) {
	var repo GHRepo
	base, err := GetAPIURL()
	if err != nil {
		return repo, nil
	}

	resp, err := http.Get(base)
	if err != nil || resp.StatusCode != http.StatusOK {
		return repo, err
	}
	return repo, nil
}

// CreateNewIssue creates a new issue on github
func CreateNewIssue(config FocusData, body map[string]string) error {
	owner, repo, err := GetRepOwnerAndName()
	if err != nil {
		return err
	}

	if body["title"] == "" || body["body"] == "" {
		return errors.New("cannot create issue without title or body")
	}

	title := strings.Trim(body["title"], " \n")
	issueBody := strings.Trim(body["body"], " \n")
	issueReq := github.IssueRequest{
		Title: &title,
		Body:  &issueBody,
	}

	if body["labels"] != "" {
		if strings.Contains(body["labels"], ",") {
			labels := strings.Split(body["labels"], ",")
			issueReq.Labels = &labels
		} else {
			labels := []string{strings.Trim(body["labels"], " \n")}
			issueReq.Labels = &labels
		}
	}

	username, password, err := promptCredentials()
	if err != nil {
		return err
	}

	issueReq.Assignees = &[]string{owner}
	// fmt.Println("t:", title, "B:", issueBody, issueReq.GetLabels(), issueReq.GetAssignees())

	client := github.NewClient(&http.Client{Transport: &github.BasicAuthTransport{
		Username: username,
		Password: password,
	}})
	_, _, err = client.Issues.Create(context.Background(), owner, repo, &issueReq)
	if err != nil {
		return err
	}
	return nil
}

func promptCredentials() (string, string, error) {
	var username string
	uprompt := &survey.Input{
		Message: "username",
	}
	survey.AskOne(uprompt, &username)
	var password string
	pprompt := &survey.Password{
		Message: "password",
	}
	survey.AskOne(pprompt, &password)

	if username == "" || password == "" {
		return "", "", errors.New("username or password needs to be pressent")
	}
	return username, password, nil
}

// CloseIssue is called when done command is ran on cli
func CloseIssue(issueNum int) error {
	owner, repo, err := GetRepOwnerAndName()
	if err != nil {
		return err
	}

	username, password, err := promptCredentials()
	if err != nil {
		return err
	}

	client := github.NewClient(&http.Client{Transport: &github.BasicAuthTransport{
		Username: username,
		Password: password,
	}})

	req := &github.IssueRequest{
		State: s("closed"),
	}

	_, _, err = client.Issues.Edit(context.Background(), owner, repo, issueNum, req)
	if err != nil {
		return err
	}
	return nil
}

func s(text string) *string {
	return &text
}
