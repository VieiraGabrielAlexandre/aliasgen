package store

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	sql *sql.DB
}

func dataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".local", "share", "aliasgen")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

func Open() (*DB, error) {
	dir, err := dataDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "aliasgen.db")
	sqlDB, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}
	db := &DB{sql: sqlDB}
	if err := db.init(); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return db, nil
}

func (d *DB) Close() error { return d.sql.Close() }

func (d *DB) init() error {
	_, err := d.sql.Exec(`
CREATE TABLE IF NOT EXISTS commands (
  command   TEXT PRIMARY KEY,
  uses      INTEGER NOT NULL DEFAULT 0,
  last_used INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS aliases (
  alias      TEXT PRIMARY KEY,
  expands_to TEXT NOT NULL,
  score      REAL NOT NULL,
  created_at INTEGER NOT NULL,
  shell      TEXT NOT NULL
);
`)
	return err
}

// UpsertCommand incrementa o uso do comando e atualiza last_used se for mais recente.
func (d *DB) UpsertCommand(cmd string, usedAt time.Time) error {
	tx, err := d.sql.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var uses int
	var last int64
	err = tx.QueryRow(`SELECT uses, last_used FROM commands WHERE command = ?`, cmd).Scan(&uses, &last)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		_, err = tx.Exec(`INSERT INTO commands(command, uses, last_used) VALUES(?,?,?)`,
			cmd, 1, usedAt.Unix())
		if err != nil {
			return err
		}
	default:
		if err != nil {
			return err
		}
		if usedAt.Unix() > last {
			_, err = tx.Exec(`UPDATE commands SET uses = ?, last_used = ? WHERE command = ?`,
				uses+1, usedAt.Unix(), cmd)
		} else {
			_, err = tx.Exec(`UPDATE commands SET uses = ? WHERE command = ?`,
				uses+1, cmd)
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

type CommandStat struct {
	Command  string
	Uses     int
	LastUsed time.Time
}

// TopCommands retorna os comandos mais usados.
func (d *DB) TopCommands(limit int) ([]CommandStat, error) {
	rows, err := d.sql.Query(`SELECT command, uses, last_used FROM commands ORDER BY uses DESC, last_used DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []CommandStat
	for rows.Next() {
		var c CommandStat
		var ts int64
		if err := rows.Scan(&c.Command, &c.Uses, &ts); err != nil {
			return nil, err
		}
		c.LastUsed = time.Unix(ts, 0)
		out = append(out, c)
	}
	return out, rows.Err()
}

func (d *DB) SaveAlias(sh, alias, expands string, score float64) error {
	_, err := d.sql.Exec(
		`INSERT OR REPLACE INTO aliases(alias, expands_to, score, created_at, shell) VALUES(?,?,?,?,?)`,
		alias, expands, score, time.Now().Unix(), sh,
	)
	return err
}

type AliasRow struct {
	Alias     string
	ExpandsTo string
	Score     float64
	Shell     string
	CreatedAt time.Time
}

func (d *DB) ListAliases() ([]AliasRow, error) {
	rows, err := d.sql.Query(`SELECT alias, expands_to, score, shell, created_at FROM aliases ORDER BY score DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AliasRow
	for rows.Next() {
		var r AliasRow
		var ts int64
		if err := rows.Scan(&r.Alias, &r.ExpandsTo, &r.Score, &r.Shell, &ts); err != nil {
			return nil, err
		}
		r.CreatedAt = time.Unix(ts, 0)
		out = append(out, r)
	}
	return out, rows.Err()
}
