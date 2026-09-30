package render

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Nios-V/wl/internal/entry"
)

var weekdays = [...]string{"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"}

func Day(w io.Writer, day time.Time, entries []entry.Entry) {
	fmt.Fprintf(w, "%s, %s\n", weekdays[day.Weekday()], day.Format("02/01"))
	if len(entries) == 0 {
		fmt.Fprintln(w, "  (no hay entradas)")
		return
	}

	width := 0
	for _, e := range entries {
		width = max(width, utf8.RuneCountInString(e.Text))
	}

	for _, e := range entries {
		var tags []string
		for _, tag := range e.Tags {
			tags = append(tags, "#"+tag)
		}
		line := fmt.Sprintf(" %s %-*s %s", e.CreatedAt.Local().Format("15:04"), width, e.Text, strings.Join(tags, " "))
		fmt.Fprintln(w, strings.TrimRight(line, " "))
	}
}
