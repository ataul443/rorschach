package main

import (
	"bufio"
	"fmt"
	"os"
)

type Rule struct {
	Event   string
	Pattern string
	Command string
}

type Rules = map[string][]*Rule

func ParseRulesFile(file *os.File) (Rules, error) {
	rules := make(map[string][]*Rule)
	lineIdx := 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ruleText := scanner.Text()
		if len(ruleText) == 0 || ruleText[0] == '#' {
			lineIdx += 1
			continue
		}

		r, err := ParseRule(ruleText, lineIdx)
		if err != nil {
			return nil, err
		}
		rules[r.Event] = append(rules[r.Event], r)
		lineIdx += 1
	}

	return rules, nil
}

func ParseRule(ruleText string, lineIdx int) (*Rule, error) {

	splitTokens := tokeniezer(ruleText)
	tokens := splitTokens()

	nt := len(tokens)
	if nt != 3 {
		return nil,
			fmt.Errorf("Line: %d, 3 tokens (words) required in single rule, tokens have %d !",
				lineIdx, nt)
	}

	allowedEvents := [3]string{"CREATE", "MODIFY", "DELETE"}
	eventValid := false
	for _, e := range allowedEvents {
		if e == tokens[0] {
			eventValid = true
			break
		}
	}

	if !eventValid {
		return nil, fmt.Errorf("Line: %d, invalid event '%s' found !", lineIdx, tokens[0])
	}

	return &Rule{Event: tokens[0], Pattern: tokens[1], Command: tokens[2]}, nil
}

func tokeniezer(text string) func() []string {
	idx := 0

	nt := len(text)

	getNextTokenIndex := func(text string, idx int) (string, int) {
		s := idx
		e := idx

		nt := len(text)
		for s < nt {
			if text[s] == '"' {
				break
			}

			s += 1
		}

		e = s + 1
		for e < nt {
			if text[e] == '"' {
				break
			}

			e += 1
		}

		if e == nt && text[e-1] != '"' {
			return "", -1
		}

		return text[s+1 : e], e + 1
	}

	return func() []string {
		tokens := make([]string, 0)
		for idx < nt {
			token, nextStartIdx := getNextTokenIndex(text, idx)
			if nextStartIdx == -1 {
				return tokens
			}

			tokens = append(tokens, token)
			idx = nextStartIdx
		}
		return tokens
	}
}
