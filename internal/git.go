package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
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
	repo, err := GetRepositoryURL()
	if err != nil {
		return "", err
	}

	split := strings.Split(repo, "/")
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", split[len(split)-2], split[len(split)-1])
	if strings.HasSuffix(apiURL, ".git") {
		dotIndex := strings.LastIndex(apiURL, ".")
		apiURL = apiURL[:dotIndex]
	}
	return apiURL, nil
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
