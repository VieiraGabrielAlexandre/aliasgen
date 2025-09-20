package learn

import (
	"fmt"
	"strings"
)

func Abbrev(cmd string) string {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return ""
	}
	if parts[0] == "git" && len(parts) > 1 {
		return "g" + initials(parts[1:])
	}
	return initials(parts)
}

func initials(parts []string) string {
	var b strings.Builder
	for _, p := range parts {
		if strings.HasPrefix(p, "-") {
			continue
		}
		b.WriteByte(p[0])
	}
	return b.String()
}

func isDangerous(cmd string) bool {
	bad := []string{"rm -rf", "mkfs", "dd if=", ">:;", ":(){", "shutdown", "reboot"}
	s := strings.ToLower(cmd)
	for _, b := range bad {
		if strings.Contains(s, b) {
			return true
		}
	}
	return false
}

func disambiguate(alias string, taken map[string]struct{}) string {
	if _, ok := taken[alias]; !ok {
		return alias
	}
	for i := 1; i < 10; i++ {
		cand := fmt.Sprintf("%s%d", alias, i)
		if _, ok := taken[cand]; !ok {
			return cand
		}
	}
	return alias + "x"
}
