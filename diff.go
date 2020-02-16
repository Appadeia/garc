package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/xanzy/go-gitlab"
)

func Diff(c *cli.Context) error {
	config := GrabConfigForRepo()
	client := config.GetClient()
	namespace, proj := GetParentProjectName()

	if HasModifications() {
		result := PromptInlineChoice("There are unstaged files. Commit using... ", "All Changes", "Without Deletions", "Without New Files", "Abort")
		message := ""
		if result == "All Changes" {
			message = PromptInEditor("", "Commit message:")
			RunCommand("git", "add", "-A")
		} else if result == "Without Deletions" {
			message = PromptInEditor("", "Commit message:")
			RunCommand("git", "add", ".")
		} else if result == "Without New Files" {
			message = PromptInEditor("", "Commit message:")
			RunCommand("git", "add", "-u")
		} else {
			println("Aborting.")
			os.Exit(0)
		}
		RunCommand("git", "commit", "-m", message)
		RunCommand("git", "push", "origin", "HEAD")

		parent, _, err := client.Projects.GetProject(namespace+"/"+proj, nil)
		CheckErr(err)

		branch := CurrentBranchName()

		mrs, _, err := client.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
			SourceBranch: &branch,
		})
		CheckErr(err)

		sourceProj, err := GetProject(client, GetProjectNameString())

		hasOpenMR := false
		for _, mr := range mrs {
			if mr.SourceProjectID == sourceProj.ID {
				hasOpenMR = true
			}
		}

		if !hasOpenMR {
			title := PromptInlineAnything("Merge request title")
			desc := PromptInEditor(`Summary:

Test Plan:

`, fmt.Sprintf("Creating a merge request to %s/%s", namespace, proj))

			branches, _, err := client.Branches.ListBranches(namespace+"/"+proj, nil)
			CheckErr(err)

			var choices []string

			for _, branch := range branches {
				choices = append(choices, branch.Name)
			}

			targetBranch := PromptInlineChoice("Target branch", choices...)
			shouldSquash := true

			mr, _, err := client.MergeRequests.CreateMergeRequest(GetProjectNameString(), &gitlab.CreateMergeRequestOptions{
				Title:           &title,
				Description:     &desc,
				SourceBranch:    &branch,
				TargetBranch:    &targetBranch,
				TargetProjectID: &parent.ID,
				Squash:          &shouldSquash,
			})
			CheckErr(err)

			PrettyPrint(mr)
		}
	} else {
		println("No changes in working directory.")
	}
	return nil
}
