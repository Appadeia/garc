package main

import (
	"fmt"
	"log"

	"github.com/xanzy/go-gitlab"
)

func (remote Remote) GetClient() *gitlab.Client {
	if remote.RemoteURL == "" || remote.Token == "" {
		log.Fatal("Please configure a remote for this repository")
	}
	client := gitlab.NewClient(nil, remote.Token)
	client.SetBaseURL(fmt.Sprintf("https://%s/api/v4", remote.RemoteURL))
	return client
}

func GetProject(client *gitlab.Client, project string) (*gitlab.Project, error) {
	proj, _, err := client.Projects.GetProject(project, nil)
	return proj, err
}

func IsProjectFork(client *gitlab.Client, project string) bool {
	proj, _, err := client.Projects.GetProject(project, nil)
	CheckErr(err)
	return proj.ForkedFromProject != nil
}
