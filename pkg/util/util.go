package util

import (
	"fmt"
	"strings"
)

func StringToUnit(unit string) string {
	if len(unit) < 2 {
		return strings.ToUpper(fmt.Sprintf("%-2s", unit))
	}

	return strings.ToUpper(unit[:2])
}
