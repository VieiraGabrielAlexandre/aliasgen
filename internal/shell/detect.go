package shell

import (
	"os"
	"path/filepath"
)

// Detect retorna bash|zsh|fish|unknown
func Detect() string {
	if sh := filepath.Base(os.Getenv("SHELL")); sh != "" {
		switch sh {
		case "bash", "zsh", "fish":
			return sh
		}
	}
	return "unknown"
}
