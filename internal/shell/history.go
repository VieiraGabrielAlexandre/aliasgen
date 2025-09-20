package shell

import (
	"os"
	"path/filepath"
)

// Paths retorna caminhos comuns de histórico
func Paths() []string {
	home, _ := os.UserHomeDir()
	return []string{
		filepath.Join(home, ".bash_history"),
		filepath.Join(home, ".zsh_history"),
		filepath.Join(home, ".local/share/fish/fish_history"),
	}
}
