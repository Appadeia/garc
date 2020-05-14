package app

import (
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

var WorkingRepository *git.Repository

func init() {
	rand.Seed(time.Now().UnixNano())
	wd, err := os.Getwd()
	CheckErr(err)
	WorkingRepository, err = git.PlainOpen(wd)
	CheckErr(err)
}

func GetRemoteOriginURL() string {
	command := exec.Command("git", "remote", "get-url", "origin")
	url, err := command.CombinedOutput()
	CheckErr(err, string(url))
	return strings.TrimSpace(string(url))
}

func GetParentOriginURL() string {
	command := exec.Command("git", "remote", "get-url", "upstream")
	url, err := command.CombinedOutput()
	CheckErr(err, string(url))
	return strings.TrimSpace(string(url))
}

func GrabConfigForRepo() Remote {
	url := GetRemoteOriginURL()

	for _, remote := range Config.Remotes {
		if strings.Contains(url, remote.RemoteURL) {
			return remote
		}
	}

	return Remote{}
}

func GetProjectNameString() string {
	namespace, project := GetProjectName()
	return namespace + "/" + project
}

func CurrentBranchName() string {
	command := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branch, err := command.CombinedOutput()
	CheckErr(err, string(branch))
	return strings.TrimSpace(string(branch))
}

func _getProjectName(url string) (string, string) {
	if strings.HasPrefix(url, "git") { // Git URL! (git@invent.kde.org:cblack/hig-kde-org.git)
		split := strings.Split(url, ":")
		if len(split) < 2 {
			ErrorOutput("Malformed URL", url)
		}

		trimmed := strings.TrimSuffix(split[1], ".git")
		trimmedSplit := strings.Split(trimmed, "/")
		if len(trimmedSplit) != 2 {
			ErrorOutput("Malformed project name", trimmed)
		}

		return trimmedSplit[0], trimmedSplit[1]
	} else if strings.HasPrefix(url, "http") { // HTTP/S URL! (https://invent.kde.org/cblack/hig-kde-org.git)
		split := strings.Split(url, "/")
		if len(split) < 3 {
			ErrorOutput("Malformed URL", url)
		}

		trimmed := strings.TrimSuffix(strings.Join(split[len(split)-2:], "/"), ".git")
		trimmedSplit := strings.Split(trimmed, "/")
		if len(trimmedSplit) != 2 {
			ErrorOutput("Malformed project name", trimmed)
		}

		return trimmedSplit[0], trimmedSplit[1]
	} else {
		ErrorOutput("Unrecognized URL scheme", url)
	}
	return "", ""
}

func GetParentProjectName() (string, string) {
	return _getProjectName(GetParentOriginURL())
}

func GetProjectName() (string, string) {
	return _getProjectName(GetRemoteOriginURL())
}

func HasModifications() bool {
	wt, err := WorkingRepository.Worktree()
	CheckErr(err)
	status, err := wt.Status()
	CheckErr(err)
	for _, status := range status {
		if status.Staging != git.Unmodified {
			return true
		}
	}
	return false
}
