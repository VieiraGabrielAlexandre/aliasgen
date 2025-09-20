package learn

import (
	"regexp"
	"strings"
)

var num = regexp.MustCompile(`\b\d+\b`)
var paths = regexp.MustCompile(`/[^\s]+`)

// Normalize remove variações voláteis e colapsa espaços
func Normalize(raw string) string {
	s := strings.TrimSpace(raw)
	s = strings.Join(strings.Fields(s), " ")
	s = num.ReplaceAllString(s, "<n>")
	s = paths.ReplaceAllString(s, "<path>")
	return s
}
