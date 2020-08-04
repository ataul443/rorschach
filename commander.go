package main

import (
	"github.com/ataul443/log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	BASPATH_ENV = "BASEPATH"
	FULLPATH_ENV = "FULLPATH"
)

type Commander struct {
	stopCh  chan struct{}
	eventCh <-chan Event
	logger     log.Logger
}

func NewCommander(eventCh <-chan Event) *Commander {
	return &Commander{
		stopCh:  make(chan struct{}),
		eventCh: eventCh,
		logger:     log.Named("commander"),
	}
}

// Run will start the commander and waits for events.
func (c *Commander) Run() {
	c.logger.Infof("Starting commander ...")
	c.listen()
}

// Stop will close all running workers also stops
// listening to any new events
func (c *Commander) Stop() {
	c.logger.Infof("Got signal for shutdown.")
	c.stopCh <- struct{}{}

}

func (c *Commander) listen() {
	for {
		select {

		case e := <-c.eventCh:
			// Execute command action
			c.execute(e)
		case <-c.stopCh:
			c.logger.Infof("Stopped listening for events.")
			return
		}
	}
}

func (c *Commander) execute(event Event) {
	ecmd := event.Command()
	fullPathRe := regexp.MustCompile("\\$\\{FULLPATH\\}")
	basePathRe := regexp.MustCompile("\\$\\{BASEPATH\\}")

	expandedFullPathVarCmd := fullPathRe.ReplaceAllString(ecmd, event.AbsPath())
	finalCmd := basePathRe.ReplaceAllString(expandedFullPathVarCmd, event.BasePath())

	cmdArgs := strings.Split(finalCmd, " ")
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	// Set ENV vars
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "%s=%s", BASPATH_ENV, event.BasePath())
	cmd.Env = append(cmd.Env, "%s=%s", FULLPATH_ENV, event.AbsPath())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	go func() {
		err := cmd.Run()
		if err != nil {
			c.logger.Errorf("error executing `%s`: %s", finalCmd, err)
		}
		for {
			select {
			case <-c.stopCh:
				return
			}
		}
	}()
}
