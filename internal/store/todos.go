package store

import (
	"database/sql"
	"time"

	"github.com/Nios-V/wl/internal/entry"
)

func (s *Store) AddTodo(text string, at time.Time) (int64, error) {
	res, err := s.db.Exec(`INSERT INTO todos (created_at, text) VALUES (?, ?)`, formatTime(at), text)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) OpenTodos() ([]entry.Todo, error) {
	rows, err := s.db.Query(`SELECT id, created_at, text FROM todos WHERE done_at IS NULL ORDER BY created_at, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entry.Todo
	for rows.Next() {
		var t entry.Todo
		var created string
		if err := rows.Scan(&t.ID, &created, &t.Text); err != nil {
			return nil, err
		}
		if t.CreatedAt, err = time.Parse(time.RFC3339, created); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) CompleteTodo(id int64, at time.Time) error {
	res, err := s.db.Exec(`UPDATE todos SET done_at = ? WHERE id = ? AND done_at IS NULL`, formatTime(at), id)
	if err != nil {
		return err
	}
	return expectOne(res)
}

func (s *Store) DeleteTodo(id int64) error {
	res, err := s.db.Exec(`DELETE FROM todos WHERE id = ? AND done_at IS NULL`, id)
	if err != nil {
		return err
	}
	return expectOne(res)
}

// retorna ErrNotFound si la sentencia no afectó ninguna fila
func expectOne(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
