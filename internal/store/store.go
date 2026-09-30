package store

import (
	"database/sql"
	"errors"
	"sort"
	"strings"
	"time"

	"github.com/Nios-V/wl/internal/entry"
	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("no encontrado")

type Store struct {
	db *sql.DB
}

var migrations = [][]string{
	{
		`CREATE TABLE entries (
			id         INTEGER PRIMARY KEY,
			created_at TEXT NOT NULL,
			text       TEXT NOT NULL
		)`,
		`CREATE TABLE entry_tags (
			entry_id INTEGER NOT NULL REFERENCES entries(id) ON DELETE CASCADE,
			tag      TEXT NOT NULL,
			PRIMARY KEY (entry_id, tag)
		)`,
		`CREATE INDEX idx_entries_created_at ON entries(created_at)`,
	},
}

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, err
	}

	s := &Store{db: db}
	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	if _, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS schema_version (version INTEGER NOT NULL)`); err != nil {
		return err
	}
	var current int
	if err := s.db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_version`).Scan(&current); err != nil {
		return err
	}
	for i := current; i < len(migrations); i++ {
		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		for _, stmt := range migrations[i] {
			if _, err := tx.Exec(stmt); err != nil {
				tx.Rollback()
				return err
			}
		}
		if _, err := tx.Exec(`INSERT INTO schema_version (version) VALUES (?)`, i+1); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func formatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func (s *Store) Add(e entry.Entry) (int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`INSERT INTO entries (created_at, text) VALUES (?, ?)`, formatTime(e.CreatedAt), e.Text)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	for _, tag := range e.Tags {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO entry_tags (entry_id, tag) VALUES (?, ?)`, id, tag); err != nil {
			return 0, err
		}
	}
	return id, tx.Commit()
}

const selectEntries = `
SELECT e.id, e.created_at, e.text, COALESCE(GROUP_CONCAT(t.tag, ','), '')
FROM entries e
LEFT JOIN entry_tags t ON t.entry_id = e.id`

func (s *Store) ListBetween(from, to time.Time) ([]entry.Entry, error) {
	rows, err := s.db.Query(selectEntries+`
		WHERE e.created_at >= ? AND e.created_at < ?
		GROUP BY e.id
		ORDER BY e.created_at, e.id`, formatTime(from), formatTime(to))
	if err != nil {
		return nil, err
	}
	return scanEntries(rows)
}

func (s *Store) Last() (entry.Entry, error) {
	rows, err := s.db.Query(selectEntries + ` GROUP BY e.id ORDER BY e.id DESC LIMIT 1`)
	if err != nil {
		return entry.Entry{}, err
	}
	es, err := scanEntries(rows)
	if err != nil {
		return entry.Entry{}, err
	}
	if len(es) == 0 {
		return entry.Entry{}, ErrNotFound
	}
	return es[0], nil
}

func (s *Store) Delete(id int64) error {
	res, err := s.db.Exec(`DELETE FROM entries WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanEntries(rows *sql.Rows) ([]entry.Entry, error) {
	defer rows.Close()
	var out []entry.Entry
	for rows.Next() {
		var e entry.Entry
		var created, tags string
		if err := rows.Scan(&e.ID, &created, &e.Text, &tags); err != nil {
			return nil, err
		}
		t, err := time.Parse(time.RFC3339, created)
		if err != nil {
			return nil, err
		}
		e.CreatedAt = t
		if tags != "" {
			e.Tags = strings.Split(tags, ",")
			sort.Strings(e.Tags)
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
