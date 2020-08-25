package internal

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IssueFile is the file format used by focus to create a new issue from command line
// the syntax makes it easy to create an issue quickly from the terminal
type IssueFile struct {
}

const (
	issueFileTempl = `
// here mention the title of your git issue
@title:

// detailed description
@body:

// if you want to tag this issue to a milestone mention it below
@milestone:


// comma separated tags for your issue
@labels:

// the username from git to assign this issue
@assignee:

`
)

var allowedIssueFileKeys = []string{
	"title", "body", "milestone", "labels", "assignee",
}

func getIssueFilePath() string {
	homedir, _ := os.UserHomeDir()
	return filepath.Join(homedir, ".FocusFile")
}

// CreateNewFile creates a new FocusFile for adding a new issue
func CreateNewFile() (string, error) {
	fp := getIssueFilePath()
	if f, _ := os.Stat(fp); f != nil {
		os.Remove(fp)
	}

	err := ioutil.WriteFile(fp, []byte(issueFileTempl), 0755)
	if err != nil {
		return "", err
	}

	return fp, nil
}

// OpenIssueFile opens the issue file inside the configured editor
// or inside nano
func OpenIssueFile(path string, editorCommand string) error {
	os.Chdir(filepath.Base(path))
	cmd := exec.Command(editorCommand, path)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ParseIssueFile parses file into a map of key value pairs which can
// be sent to Github API to create an issue
func ParseIssueFile() (map[string]string, error) {
	requestObject := make(map[string]string)
	fp := getIssueFilePath()
	if f, _ := os.Stat(fp); f == nil {
		return requestObject, errors.New("create a issue file by running: focus create")
	}
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return requestObject, err
	}

	text := string(b)
	lines := strings.Split(strings.Replace(text, "\r\n", "\n", -1), "\n")
	var keyContent string
	var operatingKey string
	for _, line := range lines {
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "\n") || len(line) == 0 {
			continue
		}

		// @ is the trigger and : is ending of the key
		if strings.HasPrefix(line, "@") {
			if operatingKey != "" {
				requestObject[operatingKey] = keyContent
				keyContent = ""
			}

			replacer := strings.NewReplacer("@", "", ":", "")
			key := replacer.Replace(line)
			if includes(key, allowedIssueFileKeys) {
				operatingKey = key
			} else {
				return requestObject, fmt.Errorf("%s not allowed here", key)
			}
		} else {
			keyContent += line + "\n"
		}
	}
	// until next @ all the lines are appended to the key
	return requestObject, nil
}

func ResetFocusFile() error {
	fpath := getIssueFilePath()
	err := os.Remove(fpath)
	if err != nil {
		return err
	}

	if _, err := CreateNewFile(); err != nil {
		return err
	}

	return nil
}

func includes(target string, list []string) bool {
	for _, l := range list {
		if l == target {
			return true
		}
	}

	return false
}
