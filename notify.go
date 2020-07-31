package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

//File stores the name of the file and event associated with it
type File struct {
	Name  string
	Event string
}

var (
	//CreateEvent stores files having create event
	CreateEvent = make(chan File)
	//ModifyEvent stores files having modify event
	ModifyEvent = make(chan File)
	//DeleteEvent stores files having modify event
	DeleteEvent = make(chan File)
)

//isFileorDir checks if the given path corresponds to file or directory
//Return true if given path is directory
//Return false if given path is file
func isFileOrDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	}
	return false
}

func watchMan(watcher *fsnotify.Watcher) {
	for {

		time.Sleep(10 * time.Second)
		fmt.Println("------------")
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			} else if event.Op&fsnotify.Create == fsnotify.Create {

				CreateEvent <- File{
					Name:  event.Name,
					Event: event.Op.String(),
				}

				if isFileOrDir(event.Name) {
					watcher.Add(event.Name)
				}

			} else if event.Op&fsnotify.Write == fsnotify.Write {

				ModifyEvent <- File{
					Name:  event.Name,
					Event: event.Op.String(),
				}

			} else if event.Op&fsnotify.Remove == fsnotify.Remove {

				DeleteEvent <- File{
					Name:  event.Name,
					Event: event.Op.String(),
				}

			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}

}

//assignDirectory assigns directories, to watcher, to be watched.
func assignDirectory(wd string, watcher *fsnotify.Watcher) {

	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatalln(err)
	}
	err = watcher.Add(wd)
	for _, f := range files {
		if f.Mode().IsDir() {

			temp := wd + "/" + f.Name()
			assignDirectory(temp, watcher)

		}
		if err != nil {
			log.Fatalln(err)
		}
	}

}

//CreateNewWatcher  creates new watcher to passed directory and its subdirectories(if any).
func CreateNewWatcher(wd string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}

	defer watcher.Close()

	done := make(chan bool)

	go watchMan(watcher)

	assignDirectory(wd, watcher)

	if err != nil {
		log.Fatalln(err)
	}
	<-done

}
