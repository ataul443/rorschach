package cmd

import (
	"testing"
)

//TestParseFlagBasic is a test for ParseFlag()

func TestParseFlagBasic(t *testing.T) {
	expectedAns := Flags{
		Interval:  5,
		RulesFile: "rules",
	}
	testAns := ParseFlag()
	if !((expectedAns.Interval == testAns.Interval) && (expectedAns.RulesFile == testAns.RulesFile)) {
		t.Error("Unable to parse command line args")
	}

}
