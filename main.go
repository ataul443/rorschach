package main

import (
	"github.com/ataul443/log"
)

func main() {
	conf, err := NewConfiguration("config.toml")
	if err != nil {
		log.Errorf("could not read config.toml: %s", err)
	}

	w := NewWatcher(conf.WatchDir, conf.Rules,
		WithScanInterval(conf.ScanInterval),
		WithExcludeDir(conf.ExcludeDir),
		WithLogger(log.Named("watcher")))

	eventCh := w.Watch()
	c := NewCommander(eventCh)
	c.Run()

}
