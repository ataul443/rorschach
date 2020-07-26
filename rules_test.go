package main

import (
	"reflect"
	"testing"
)

func TestTokeniezer(t *testing.T) {
	type payload struct {
		text   string
		tokens []string
	}

	tcs := []payload{
		{
			text:   `"MODIFY"   "*.c"     "cc -o ${BASEPATH} ${FULLPATH}"`,
			tokens: []string{"MODIFY", "*.c", "cc -o ${BASEPATH} ${FULLPATH}"},
		},
		{
			text:   `"MODIFY" "cc -o ${BASEPATH} ${FULLPATH}"`,
			tokens: []string{"MODIFY", "cc -o ${BASEPATH} ${FULLPATH}"},
		},
		{
			text:   `"MODIFY" "cc -o ${BASEPATH} ${FULLPATH}`,
			tokens: []string{"MODIFY"},
		},
	}

	for _, tc := range tcs {
		splitTokenFn := tokeniezer(tc.text)
		tokens := splitTokenFn()

		if !reflect.DeepEqual(tokens, tc.tokens) {
			t.Fatalf("Got tokens: %v,\nExpected tokens: %v\n", tokens, tc.tokens)
		}
	}
}

func TestParseRule(t *testing.T) {
	type payload struct {
		text string
		rule *Rule
	}

	tcs := []payload{
		{
			text: `"MODIFY"   "*.c"     "cc -o ${BASEPATH} ${FULLPATH}"`,
			rule: &Rule{Event: "MODIFY", Pattern: "*.c", Command: "cc -o ${BASEPATH} ${FULLPATH}"},
		},
		{
			text: `"MODIFY" "cc -o ${BASEPATH} ${FULLPATH}"`,
			rule: nil,
		},
		{
			text: `"MODIFY" "cc -o ${BASEPATH} ${FULLPATH}`,
			rule: nil,
		},
	}

	for lineIdx, tc := range tcs {
		rule, _ := ParseRule(tc.text, lineIdx)
		if !reflect.DeepEqual(rule, tc.rule) {
			t.Fatalf("Got rule: %v,\nExpected rule: %v\n", rule, tc.rule)
		}
	}
}
