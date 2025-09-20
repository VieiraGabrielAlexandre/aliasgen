package aliasfmt

import "fmt"

func Bash(alias, expands string) string {
	return fmt.Sprintf("alias %s='%s'", alias, expands)
}

func Zsh(alias, expands string) string {
	return fmt.Sprintf("alias %s='%s'", alias, expands)
}

func Fish(alias, expands string) string {
	return fmt.Sprintf("alias %s '%s'", alias, expands)
}
