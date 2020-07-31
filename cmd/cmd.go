package cmd

import (
	"flag"
	"log"
	"os"
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

//FetchDirectory fetch directory name to be monitored from the command line argument
//It returns string value which is a fullpath to directory to be monitored
func FetchDirectory() string {

	args := os.Args
	size := len(args)
	if size == 1 {
		log.Fatalln("Please mention directory/file to be monitored")
	}
	//last element from of the arguments slice
	basepath := args[size-1]
	//current directory
	fullpath, _ := os.Getwd()

	if basepath != "." {
		fullpath += basepath
	}
	//check if directory/file exists or not
	if _, err := os.Stat(fullpath); os.IsNotExist(err) {
		log.Fatalln("Directory Doesn't exists")
	}

	return fullpath
}
