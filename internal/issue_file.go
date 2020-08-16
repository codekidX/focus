package internal

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// IssueFile is the file format used by focus to create a new issue from command line
// the syntax makes it easy to create an issue quickly from the terminal
type IssueFile struct {
}

const issueFileTempl = `
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

// CreateNewFile creates a new FocusFile for adding a new issue
func CreateNewFile() (string, error) {
	homedir, _ := os.UserHomeDir()
	fp := filepath.Join(homedir, ".FocusFile")
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
