package strings

import "strings"

// Split Fix If s does not contain sep and sep is not empty, Split returns a slice of length 1 whose only element is s.
func Split(s, sep string) (results []string) {
	results = strings.Split(s, sep)
	if len(results) == 1 {
		if results[0] == "" {
			results = []string{}
		}
	}
	return
}
