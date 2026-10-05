package store

import (
	"database/sql"
	"errors"
	"time"
)

const (
	KindDay    = "day"
	KindSprint = "sprint"
)

func (s *Store) SetSummary(kind, period, text string, at time.Time) error {
	_, err := s.db.Exec(`
		INSERT INTO summaries (kind, period, text, updated_at) VALUES (?, ?, ?, ?)
		ON CONFLICT (kind, period) DO UPDATE SET text = excluded.text, updated_at = excluded.updated_at`,
		kind, period, text, formatTime(at))
	return err
}

// Summary devuelve "" si no hay resumen para ese periodo
func (s *Store) Summary(kind, period string) (string, error) {
	var text string
	err := s.db.QueryRow(`SELECT text FROM summaries WHERE kind = ? AND period = ?`, kind, period).Scan(&text)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return text, err
}
