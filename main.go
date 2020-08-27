package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

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
page - focus on page ${number}
open - open issue in browser

create - create issue from your local machine
push - push created issue to GitHub
reset - cleans the previous issue file that you wrote and replaces with fresh issue file
status - shows the pending issue(if any) which is ready to push

done - mark an issue as closed/done
todo - create a todo for your current repository
todos - list TODOs for your current repository
todos [rm] $id - removes a TODO for given $id from this repository
`

// TODO: add this to help command after implementing this freature

func main() {
	flag.Parse()

	args := flag.Args()
	fd, err := internal.GetFocusData()
	errp(err)

	if len(args) == 0 {
		issues, err := internal.ListIssues("")
		if err != nil {
			t.Println("cannot fetch issues for this directory, is it a git repository?",
				tint.Red.Bold())
			return
		}
		if len(issues) == 0 {
			errp(errors.New("no issues in this repository, do: focus create"))
		}
		displayIssueList(issues)
	} else if len(args) > 0 {
		command := args[0]
		err := checkCommand(command, args[1:], fd)
		errp(err)
	}
}

func errp(err error) {
	if err != nil {
		// t.Println(err.Error(), tint.Red.Bold())
		// os.Exit(1)
		panic(err)
	}
}

func checkCommand(command string, args []string, fd internal.FocusData) error {
	switch command {
	case "create":
		path, err := internal.CreateNewFile()
		if err != nil {
			return err
		}
		return internal.OpenIssueFile(path, fd.Editor)
	case "push":
		body, err := internal.ParseIssueFile()
		if err != nil {
			return err
		}

		// fmt.Printf("%T %v >>\n", body["labels"], body)

		err = internal.CreateNewIssue(fd, body)
		if err != nil {
			return err
		}

		msg := t.Exp(fmt.Sprintf("@(%s) is pushed to GitHub", strings.Trim(body["title"], " \n")),
			tint.Yellow)
		fmt.Println(msg)
		return nil
	case "reset":
		return internal.ResetFocusFile()
	case "status":
		ff, err := internal.ParseIssueFile()
		if err != nil {
			return err
		}

		if ff["title"] != "" {
			title := strings.Trim(ff["title"], " \n")
			exp := t.Exp(fmt.Sprintf("(@(%s)) -> @(ready to push)", title),
				tint.Green, tint.Yellow.Bold())
			fmt.Println(exp)
		} else {
			return errors.New("no pending issue file to push")
		}
		return nil
	case "done":
		if len(args) == 0 {
			msg := "done command requires a issue number for marking it as closed, run: focus"
			return errors.New(msg)
		}

		id, err := strconv.Atoi(args[0])
		// TODO: check if not number then check if the arg[0] is "todo"
		// if it is todo then mark todo at index for this repo
		// as done -- CALL RemoveTODO(:index)
		if err != nil {
			return errors.New("not a id: number/int")
		}
		err = internal.CloseIssue(id)
		if err != nil {
			return err
		}
		fmt.Println(t.Exp(fmt.Sprintf("@(Closed) issue %d", id), tint.Red))
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

		// TODO: pass this as query params not as a string
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
	case "todo":
		if len(args) == 0 {
			msg := "enter text to save todo -> example: focus todo \"we need to do this\""
			return errors.New(msg)
		}

		if len(args) > 1 {
			return errors.New("you need to pass todo in quotes for it to identify as a string")
		}

		err := internal.SaveTODO(fd, args[0])
		if err != nil {
			return err
		}
		return nil
	case "todos":
		if len(args) > 0 && args[0] == "rm" {
			if len(args) > 1 {
				id, err := strconv.Atoi(args[1])
				if err != nil {
					return errors.New("not a id: number/int")
				}
				// we remove todo at given index(human)
				err = internal.RemoveTODO(fd, id)
				if err != nil {
					return err
				}
			} else {
				// output that rm requires an index
				return errors.New("rm requires an argument for index of TODO to remove")
			}
		}
		internal.ListTODOs(fd)
		return nil
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
