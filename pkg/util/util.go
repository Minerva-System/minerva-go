package util

import (
	"fmt"
	"strings"
	"regexp"
	"errors"
)

func StringToUnit(unit string) string {
	if len(unit) < 2 {
		return strings.ToUpper(fmt.Sprintf("%-2s", unit))
	}

	return strings.ToUpper(unit[:2])
}

func HygienizeSlug(slug string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(slug))
	if regexp.MustCompile(`\s*`).MatchString(s) {
		return slug, errors.New("Slug must not have whitespaces")
	}
	return s, nil
}
