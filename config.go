package main

import (
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

func LoadConfiguration() Configuration {
	home, err := os.UserHomeDir()
	CheckErr(err)
	configPath := path.Join(home, ".garcrc")

	var config Configuration
	_, err = toml.DecodeFile(configPath, &config)
	CheckErr(err)

	return config
}
