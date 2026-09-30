package entry

import "time"

type Entry struct {
	ID        int64
	CreatedAt time.Time
	Text      string
	Tags      []string
}
