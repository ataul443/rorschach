package main

import (
	"github.com/BurntSushi/toml"
)

type Configuration struct {
	WatchDir     string   `toml:"watchDir"`
	ExcludeDir   []string `toml:"excludeDir"`
	ScanInterval int      `toml:"scanInterval"`
	Rules        []Rule   `toml:"rule"`
}

func NewConfiguration(configPath string) (*Configuration, error) {
	conf := new(Configuration)
	_, err := toml.DecodeFile(configPath, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
