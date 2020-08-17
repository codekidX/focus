package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/codekidX/focus/internal"
	"github.com/pkg/browser"
	"github.com/printzero/tint"
)

var t = tint.Init()

var helpText = `
🎯 @(Focus!)
Run command with no arguments to list issues of your current repository

@(COMMANDS:)
on - focus on issue with ID
open - open issue in browser
create - create issue from your local machine
push - push created issue using the 'create' command
`

func main() {
	flag.Parse()

	args := flag.Args()
	config, err := internal.GetConfig()
	errp(err)

	if len(args) == 0 {
		issues, err := internal.ListIssues("")
		if err != nil {
			errp(err)
		}
		displayIssueList(issues)
	} else if len(args) > 0 {
		command := args[0]
		err := checkCommand(command, args[1:], config)
		errp(err)
	}

}

func errp(err error) {
	if err != nil {
		t.Println(err.Error(), tint.Red.Bold())
		os.Exit(1)
	}
}

func checkCommand(command string, args []string, config internal.Config) error {
	switch command {
	case "create":
		path, err := internal.CreateNewFile()
		if err != nil {
			return err
		}
		return internal.OpenIssueFile(path, config.Editor)
	case "push":
		body, err := internal.ParseIssueFile()
		if err != nil {
			return err
		}

		err = internal.CreateNewIssue(config, body)
		if err != nil {
			return err
		}
		return nil
	case "on":
		if len(args) == 0 {
			msg := "show command requires the issue id to be passed as argument, for more info " +
				"do: focus list"
			return errors.New(msg)
		}

		id, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("not a id: number/int")
		}

		issue, err := internal.GetIssue(id)
		if err != nil {
			return err
		}
		displayFullIssue(issue)
		return nil
	case "page":
		if len(args) == 0 {
			msg := "page command requires the page number as argument, for more info " +
				"run: focus"
			return errors.New(msg)
		}
		page, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("not a valid page number: number/int")
		}
		issues, err := internal.ListIssues(fmt.Sprintf("?page=%d", page))
		if err != nil {
			return err
		}
		displayIssueList(issues)
		return nil
	case "open":
		if len(args) == 0 {
			msg := "open command requires the issue id to be passed as argument, for more info " +
				"run: focus"
			return errors.New(msg)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("not a id: number/int")
		}

		issue, err := internal.GetIssue(id)
		if err != nil {
			return err
		}

		return browser.OpenURL(issue.HTMLURL)
	case "h", "help", "?":
		out := t.Exp(helpText, tint.Cyan, tint.Yellow)
		fmt.Println(out)
		return nil
	default:
		fmt.Println("no such command")
		return nil
	}
}

func getTitleTintExpression(issue internal.GHIssue) string {
	msg := fmt.Sprintf("@(#%d) %s", issue.Number, issue.Title)
	return t.Exp(msg, tint.Green)
}

func displayIssueList(issues []internal.GHIssue) {
	for _, issue := range issues {
		fmt.Println(getTitleTintExpression(issue))
	}
}

func displayFullIssue(issue internal.GHIssue) {
	titleExp := getTitleTintExpression(issue)
	fmt.Println("\n" + titleExp + "\n")
	displayData("Status:", issue.State)
	displayData("Body:", "\n"+issue.Body+"\n")
}

func displayData(name string, value string) {
	exp := t.Exp(fmt.Sprintf("@(%s) %s", name, value), tint.Yellow)
	fmt.Println(exp)
}
