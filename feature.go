package main

import (
	"github.com/urfave/cli/v2"
	"github.com/xanzy/go-gitlab"
)

func Feature(c *cli.Context) error {
	config := GrabConfigForRepo()
	client := config.GetClient()
	_, project := GetProjectName()

	if c.Args().First() == "" {
		ErrorOutput("Please provide a feature repo name")
	}
	if IsProjectFork(client, GetProjectNameString()) {
		RunCommand("git", "checkout", "master")

		RunCommand("git", "pull")

		RunCommand("git", "checkout", "-b", c.Args().First())
	} else {
		if fork, err := GetProject(client, config.Username+"/"+project); err == nil {
			url := GetRemoteOriginURL()

			RunCommand("git", "remote", "set-url", "origin", fork.SSHURLToRepo)

			RunCommand("git", "remote", "add", "upstream", url)

			RunCommand("git", "pull")

			RunCommand("git", "checkout", "-b", c.Args().First())
		} else {
			fork, _, err := client.Projects.ForkProject(GetProjectNameString(), &gitlab.ForkProjectOptions{
				Namespace: &config.Username,
			})
			CheckErr(err)

			url := GetRemoteOriginURL()

			RunCommand("git", "remote", "set-url", "origin", fork.SSHURLToRepo)

			RunCommand("git", "remote", "add", "upstream", url)

			RunCommand("git", "pull")

			RunCommand("git", "checkout", "-b", c.Args().First())
		}
	}
	return nil
}
