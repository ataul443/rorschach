package main

import "testing"

func TestMatchWildcard(t *testing.T) {
	type payload struct {
		pattern string
		text    string
		match   bool
	}

	tcs := []payload{
		{pattern: "*a*b", text: "adceb", match: true},
		{pattern: "*b", text: "adceb", match: true},
		{pattern: "a*c?b", text: "aacdcb", match: false},
		{pattern: "*ab***ba**b*b*aaab*b", text: "aaabababaaabaababbbaaaabbbbbbabbbbabbbabbaabbababab", match: true},
	}

	for _, tc := range tcs {
		match := MatchWildcard(tc.pattern, tc.text)
		if match != tc.match {
			t.Fatalf("text = `%s`, pattern = `%s`, got match = %v, expected match = %v\n",
				tc.text, tc.pattern, match, tc.match)
		}
	}
}

func BenchmarkMatchWildcard(b *testing.B) {
	b.ReportAllocs()
	pattern := "*ab***ba**b*b*aaab*b"
	text := "aaabababaaabaababbbaaaabbbbbbabbbbabbbabbaabbababab"
	for n := 0; n < b.N; n++ {
		MatchWildcard(pattern, text)
	}

}
