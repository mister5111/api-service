package customvalidator

import (
	"regexp"
)

func AliasValidator(alias string) bool {
	if regexp.MustCompile(`(?i)^(http://|https://|http:/|http:|http|https:/|https:|https)`).MatchString(alias) {
		return false
	}

	match, _ := regexp.MatchString(`^[a-zA-Z0-9_.-]+$`, alias)
	return match
}
