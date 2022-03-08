package utils

import "strings"

// BadSubject will do quick test on whether a subject is acceptable.
// Spaces are not allowed and all tokens should be > 0 in len.
func BadSubject(subj string) bool {
	if strings.ContainsAny(subj, " \t\r\n") {
		return true
	}

	tokens := strings.Split(subj, ".")

	for _, t := range tokens {
		if len(t) == 0 {
			return true
		}
	}

	return false
}