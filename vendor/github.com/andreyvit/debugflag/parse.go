package debugflag

import (
	"regexp"
)

const (
	Neg1 byte = '!'
	Neg2 byte = '-'
)

var separators = regexp.MustCompile("[\\s,]+")

func parse(spec string) ([]string, []string) {
	patterns := separators.Split(spec, -1)
	var enabled, disabled []string
	for _, pattern := range patterns {
		if len(pattern) == 0 {
			continue
		}
		if pattern[0] == Neg1 || pattern[0] == Neg2 {
			pattern = pattern[1:]
			if len(pattern) > 0 {
				disabled = append(disabled, pattern)
			}
		} else {
			enabled = append(enabled, pattern)
		}
	}
	return enabled, disabled
}
