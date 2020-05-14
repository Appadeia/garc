package app

import (
	"log"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type Remote struct {
	Token     string `toml:"token"`
	RemoteURL string `toml:"url"`
	Username  string `toml:"username"`
}

type Configuration struct {
	Remotes []Remote `toml:"remote"`
}

var Config Configuration

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	configPath := path.Join(home, ".garcrc")

	_, err = toml.DecodeFile(configPath, &Config)
	if err != nil {
		log.Fatal(err)
	}
}
