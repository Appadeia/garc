package main

import (
	"os"

	"github.com/pontaoski/garc/app"
	"github.com/urfave/cli"
	"github.com/xanzy/go-gitlab"
)

func Main(c *cli.Context) error {
	config := app.GrabConfigForRepo()
	client := config.GetClient()
	_, project := app.GetProjectName()

	if c.Args().First() == "" {
		app.ErrorOutput("Please provide a feature repo name")
	}
	if app.IsProjectFork(client, app.GetProjectNameString()) {
		app.RunCommand("git", "checkout", "master")

		app.RunCommand("git", "pull")

		app.RunCommand("git", "checkout", "-b", c.Args().First())
	} else {
		if fork, err := app.GetProject(client, config.Username+"/"+project); err == nil {
			url := app.GetRemoteOriginURL()

			app.RunCommand("git", "remote", "set-url", "origin", fork.SSHURLToRepo)

			app.RunCommand("git", "remote", "add", "upstream", url)

			app.RunCommand("git", "pull")

			app.RunCommand("git", "checkout", "-b", c.Args().First())
		} else {
			fork, _, err := client.Projects.ForkProject(app.GetProjectNameString(), &gitlab.ForkProjectOptions{
				Namespace: &config.Username,
			})
			app.CheckErr(err)

			url := app.GetRemoteOriginURL()

			app.RunCommand("git", "remote", "set-url", "origin", fork.SSHURLToRepo)

			app.RunCommand("git", "remote", "add", "upstream", url)

			app.RunCommand("git", "pull")

			app.RunCommand("git", "checkout", "-b", c.Args().First())
		}
	}
	return nil
}

func main() {
	cliApp := cli.App{
		Action: Main,
	}
	cliApp.Run(os.Args)
}
