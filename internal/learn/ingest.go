package learn

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/VieiraGabrielAlexandre/aliasgen/internal/shell"
	"github.com/VieiraGabrielAlexandre/aliasgen/internal/store"
)

// Run faz a ingestão dos históricos conhecidos e grava no banco.
func Run() (int, error) {
	db, err := store.Open()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	paths := shell.Paths()
	count := 0
	for _, p := range paths {
		n, _ := ingestFile(db, p)
		count += n
	}
	return count, nil
}

func ingestFile(db *store.DB, path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	var n int
	switch filepath.Base(path) {
	case "fish_history":
		for sc.Scan() {
			line := sc.Text()
			// fish: "- cmd: <comando>" e, em linha posterior, "  when: <unix>"
			if strings.HasPrefix(line, "- cmd: ") {
				cmd := strings.TrimSpace(strings.TrimPrefix(line, "- cmd: "))
				cmd = Normalize(cmd)
				if cmd == "" {
					continue
				}
				ts := time.Now()
				if sc.Scan() {
					next := strings.TrimSpace(sc.Text())
					if strings.HasPrefix(next, "when: ") {
						if u, perr := parseUnix(strings.TrimSpace(strings.TrimPrefix(next, "when: "))); perr == nil {
							ts = time.Unix(u, 0)
						}
					}
				}
				if err := db.UpsertCommand(cmd, ts); err == nil {
					n++
				}
			}
		}
	default:
		// bash: linha é o comando
		// zsh extended: ": 1700000000:0;git status"
		for sc.Scan() {
			line := sc.Text()
			cmd, ts := parseBashZsh(line)
			cmd = Normalize(cmd)
			if cmd == "" {
				continue
			}
			if err := db.UpsertCommand(cmd, ts); err == nil {
				n++
			}
		}
	}
	return n, sc.Err()
}

func parseUnix(s string) (int64, error) {
	s = strings.TrimSpace(s)
	// pega apenas prefixo numérico
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == 0 {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(s[:i], 10, 64)
}

func parseBashZsh(line string) (cmd string, ts time.Time) {
	ts = time.Now()
	line = strings.TrimSpace(line)
	if line == "" {
		return "", ts
	}
	// zsh formato: ": <unix>:<duration>;<cmd>"
	if strings.HasPrefix(line, ": ") {
		semi := strings.Index(line, ";")
		if semi > 0 && semi < len(line)-1 {
			header := line[:semi]
			payload := strings.TrimSpace(line[semi+1:])
			parts := strings.Split(header, ":")
			if len(parts) >= 3 {
				unixStr := strings.TrimSpace(parts[1])
				if u, err := parseUnix(unixStr); err == nil {
					ts = time.Unix(u, 0)
				}
			}
			return payload, ts
		}
	}
	// bash simples
	return line, ts
}
