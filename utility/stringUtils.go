package utility

import (
	"strings"
)

//CleanInput removes spaces from string
func CleanInput(in string) string {
	return strings.TrimSpace(in)
}
