package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"gopkg.in/src-d/go-git.v4"
)

var WorkingRepository *git.Repository

func main() {
	app := &cli.App{
		Name:  "garc",
		Usage: "Work with GitLab like it's Phabricator",
		Commands: []*cli.Command{
			{
				Name:    "feature",
				Aliases: []string{"branch", "b", "f"},
				Usage:   "Create a new feature branch",
				Action:  Feature,
			},
			{
				Name:    "diff",
				Aliases: []string{"d"},
				Usage:   "Submit or update a merge request",
				Action:  Diff,
			},
			{
				Name:  "support",
				Usage: "Get support for garc",
				Action: func(c *cli.Context) error {
					println(`If you need support, contact me at one of these places:
- @appadeia:matrix.org
- @pontaoski on telegram
- pontaoski#8758 on discord
- uhhadd AT gmail DOT com

Please note that I only support these distros
due to my knowledge of them and their quirks:

- Arch Linux
- openSUSE Leap 15.1
- openSUSE Tumblweed
- Fedora 32
- Mageia
- OpenMandriva
- Ubuntu 18.04
- Alpine Linux (or pmOS, but why are you running this on a phone?)

I do not have the time to learn about the
quirks or know where to find information about all
distros, which is why I only support these.

¿Habla español? Puede hablar conmigo en español.`)
					return nil
				},
			},
		},
	}
	rand.Seed(time.Now().UnixNano())
	wd, err := os.Getwd()
	CheckErr(err)
	WorkingRepository, err = git.PlainOpen(wd)
	CheckErr(err)

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err.Error())
	}
}
