package main

import (
	"fmt"
	"syscall"
	"time"
)

func shouldDirExclude(dir string, excludeDirs []string) bool {
	for _, v := range excludeDirs {
		if v == dir {
			return true
		}
	}
	return false
}

func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}

func searchPattern(fileInfo fileStat, event string, rules []Rule) (Rule, error) {
	for _, rule := range rules {
		if rule.Event == event {
			isFullPathMatched := MatchWildcard(rule.Pattern, fileInfo.absPath)
			isBasePathMatched := MatchWildcard(rule.Pattern, fileInfo.baseName)

			if isBasePathMatched || isFullPathMatched {
				return rule, nil
			}
		}
	}

	return Rule{}, fmt.Errorf("no rule found for file %s", fileInfo.absPath)
}
