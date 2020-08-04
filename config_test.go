package main

import "testing"

func TestConfiguration(t *testing.T) {
	configPath := "config.toml"
	c, err := NewConfiguration(configPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%v\n", *c)
}
