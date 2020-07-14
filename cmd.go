package main

import (
	"flag"
)

//Flags type contains the command line flag values
type Flags struct {
	Interval  int    //time interval between each scans (in seconds)
	RulesFile string //file name to load rules
}

//ParseFlag parse command line argument and return pointer to the struct type.
func ParseFlag() *Flags {

	interval := flag.Int("t", 5, "Time between scans (default is 5 seconds)")
	rulesfile := flag.String("f", "rules", "Load rules from this file (default is 'rules')")

	flag.Parse()

	return &Flags{
		Interval:  *interval,
		RulesFile: *rulesfile,
	}
}
