package main

type Rule struct {
	Event   string `toml:"event"`
	Pattern string `toml:"pattern"`
	Command string `toml:"command"`
}
