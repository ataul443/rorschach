package main

import (
	"os"
	"path"
	"sync"
	"syscall"
	"time"

	"github.com/ataul443/log"
)

const (
	CREATE_EVENT = "CREATE"
	MODIFY_EVENT = "MODIFY"
	DELETE_EVENT = "DELETE"
)

type fileStat struct {
	scanID           uint
	absPath          string
	baseName         string
	creationTime     time.Time
	modificationTime time.Time
}

type watcherOptn func(*Watcher)

type Watcher struct {
	scanInterval time.Duration
	rootPath     string
	logger       log.Logger
	excludeDir   []string
	rules        []Rule
	stopCh       chan struct{}
	fileStats    map[string]fileStat
	scanID       uint
}

func WithLogger(logger log.Logger) watcherOptn {
	return func(w *Watcher) {
		w.logger = logger
	}
}

func WithExcludeDir(excludeDir []string) watcherOptn {
	return func(w *Watcher) {
		w.excludeDir = excludeDir
	}
}

func WithScanInterval(timeInterval int) watcherOptn {
	return func(w *Watcher) {
		w.scanInterval = time.Second * time.Duration(timeInterval)
	}
}

func NewWatcher(rootPath string, rules []Rule, optns ...watcherOptn) *Watcher {
	w := new(Watcher)
	w.rootPath = rootPath
	w.rules = rules
	w.fileStats = make(map[string]fileStat)
	w.scanID = 0

	for _, optn := range optns {
		optn(w)
	}

	return w
}

func (w *Watcher) Stop() {
	// It might give panic if already closed
	w.stopCh <- struct{}{}
}

func (w *Watcher) Watch() chan Event {
	eventCh := make(chan Event, 10)
	// ticker := time.NewTicker(w.scanInterval)
	go func() {
		w.logger.Infof("Started watching `%s` ...", w.rootPath)
		ticker := time.NewTicker(w.scanInterval)
		for {
			select {
			case <-ticker.C:
				// Now its time for scan
				w.scan(w.rootPath, eventCh)
				w.scanID += 1
			case <-w.stopCh:
				return
			}
		}

	}()
	return eventCh
}

func (w *Watcher) scan(fullDirPath string, eventCh chan<- Event) {
	fileInfoCollector := func() chan fileStat {
		ch := make(chan fileStat)
		go func() {
			defer close(ch)
			w.extractFileInfoFromDir(fullDirPath, ch)
		}()

		return ch
	}
	fileCh := fileInfoCollector()
	w.detectEvents(fileCh, eventCh)
}

func (w *Watcher) eventDispatcher(f fileStat, eventName string, rules []Rule, eventCh chan<- Event) {
	rule, err := searchPattern(f, eventName, w.rules)
	if err != nil {
		w.logger.Debugf("%s", err)
	} else {
		w.logger.Infof(`Matched pattern "%s" on "%s"`, rule.Pattern, f.baseName)
		e := event{
			fullFilePath: f.absPath,
			baseFilePath: f.baseName,
			rule:         rule,
		}
		eventCh <- e
	}
}

func (w *Watcher) detectEvents(fileCh chan fileStat, eventCh chan<- Event) {
	for {
		select {
		case f, ok := <-fileCh:
			if !ok {
				// All scan done
				// Check for deleted files
				for k, v := range w.fileStats {
					if v.scanID < w.scanID {
						w.logger.Infof(`Detected "%s" on "%s"`, DELETE_EVENT, v.baseName)
						// "DELETE" event detected
						w.eventDispatcher(v, DELETE_EVENT, w.rules, eventCh)
						delete(w.fileStats, k)
					}
				}

				return
			}

			if w.scanID != uint(0) {
				// Check if it is an event
				oldFileInfo, ok := w.fileStats[f.absPath]
				if !ok {
					// New file found, "CREATE" event
					w.logger.Infof(`Detected "%s" on "%s"`, CREATE_EVENT, f.baseName)
					w.eventDispatcher(f, CREATE_EVENT, w.rules, eventCh)
				} else {
					isModified := oldFileInfo.modificationTime.Before(f.modificationTime)
					if isModified {
						// File modified, do something
						w.logger.Infof(`Detected "%s" on "%s"`, MODIFY_EVENT, f.baseName)
						w.eventDispatcher(f, MODIFY_EVENT, w.rules, eventCh)
					}
				}

			}
			w.fileStats[f.absPath] = f
		}
	}
}

func (w *Watcher) extractFileInfoFromDir(fullDirPath string, fileInfoCollector chan fileStat) {
	d, err := os.Open(fullDirPath)
	if err != nil {
		log.Errorf("Failed to open `%s`: %s", fullDirPath, err)
		return
	}

	// Be careful, this directory might contains millions of files
	nodesInfo, err := d.Readdir(-1)
	if err != nil {
		log.Errorf("Failed to scan `%s`: %s", fullDirPath, err)
		return
	}

	wg := &sync.WaitGroup{}
	for _, n := range nodesInfo {
		nextNode := path.Join(fullDirPath, n.Name())
		if n.IsDir() {
			// Check if directory needed to be excluded
			if shouldDirExclude(n.Name(), w.excludeDir) {
				log.Debugf("excluding %s", nextNode)
				continue
			}

			// Node is a dir, do something about it
			wg.Add(1)
			go func() {
				defer wg.Done()
				w.extractFileInfoFromDir(nextNode, fileInfoCollector)
			}()
		} else {
			// Node is a file, dispatch it to fileInfoCollector
			stat_t := n.Sys().(*syscall.Stat_t)

			fileInfoCollector <- fileStat{
				scanID:           w.scanID,
				baseName:         n.Name(),
				absPath:          nextNode,
				creationTime:     timespecToTime(stat_t.Ctim),
				modificationTime: timespecToTime(stat_t.Mtim),
			}
		}
	}
	wg.Wait()
}
