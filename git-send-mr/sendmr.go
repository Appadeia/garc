package main

import (
	"fmt"
	"os"

	"github.com/pontaoski/garc/app"
	"github.com/urfave/cli/v2"
	"github.com/xanzy/go-gitlab"
)

func Diff(c *cli.Context) error {
	config := app.GrabConfigForRepo()
	client := config.GetClient()
	namespace, proj := app.GetParentProjectName()

	if app.HasModifications() {
		result := app.PromptInlineChoice("There are unstaged files. Commit using... ", "All Changes", "Without Deletions", "Without New Files", "Abort")
		message := ""
		if result == "All Changes" {
			message = app.PromptInEditor("", "Commit message:")
			app.RunCommand("git", "add", "-A")
		} else if result == "Without Deletions" {
			message = app.PromptInEditor("", "Commit message:")
			app.RunCommand("git", "add", ".")
		} else if result == "Without New Files" {
			message = app.PromptInEditor("", "Commit message:")
			app.RunCommand("git", "add", "-u")
		} else {
			println("Aborting.")
			os.Exit(0)
		}
		app.RunCommand("git", "commit", "-m", message)
		app.RunCommand("git", "push", "origin", "HEAD")

		parent, _, err := client.Projects.GetProject(namespace+"/"+proj, nil)
		app.CheckErr(err)

		branch := app.CurrentBranchName()

		mrs, _, err := client.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
			SourceBranch: &branch,
		})
		app.CheckErr(err)

		sourceProj, err := app.GetProject(client, app.GetProjectNameString())

		hasOpenMR := false
		for _, mr := range mrs {
			if mr.SourceProjectID == sourceProj.ID {
				hasOpenMR = true
			}
		}

		if !hasOpenMR {
			title := app.PromptInlineAnything("Merge request title")
			desc := app.PromptInEditor(`Summary:

Test Plan:

`, fmt.Sprintf("Creating a merge request to %s/%s", namespace, proj))

			branches, _, err := client.Branches.ListBranches(namespace+"/"+proj, nil)
			app.CheckErr(err)

			var choices []string

			for _, branch := range branches {
				choices = append(choices, branch.Name)
			}

			targetBranch := app.PromptInlineChoice("Target branch", choices...)
			shouldSquash := true

			mr, _, err := client.MergeRequests.CreateMergeRequest(app.GetProjectNameString(), &gitlab.CreateMergeRequestOptions{
				Title:           &title,
				Description:     &desc,
				SourceBranch:    &branch,
				TargetBranch:    &targetBranch,
				TargetProjectID: &parent.ID,
				Squash:          &shouldSquash,
			})
			app.CheckErr(err)

			app.PrettyPrint(mr)
		}
	} else {
		println("No changes in working directory.")
	}
	return nil
}

func main() {
	cliApp := cli.App{
		Action: Diff,
	}
	cliApp.Run(os.Args)
}
