package entry

import "time"

type Todo struct {
	ID        int64
	CreatedAt time.Time
	Text      string
}
