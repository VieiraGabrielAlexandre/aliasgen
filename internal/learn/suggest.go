package learn

import "strings"

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
