package render

import (
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Nios-V/wl/internal/entry"
	"github.com/fatih/color"
)

var weekdays = [...]string{"Domingo", "Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado"}

var (
	summaryStyle = color.New(color.FgCyan)
	todoStyle    = color.New(color.FgYellow)
)

func Day(w io.Writer, day time.Time, entries []entry.Entry, summary string) {
	header(w, day)
	if len(entries) == 0 {
		fmt.Fprintln(w, "  (no hay entradas)")
	} else {
		writeEntries(w, entries)
	}
	writeSummary(w, summary)
}

func Sprint(w io.Writer, start, end time.Time, entries []entry.Entry, summary string) {
	fmt.Fprintf(w, "Sprint %s → %s\n", start.Format("02/01"), end.AddDate(0, 0, -1).Format("02/01"))
	if len(entries) == 0 {
		fmt.Fprintln(w, "  (no hay entradas)")
	}
	for i := 0; i < len(entries); {
		day := entry.StartOfDay(entries[i].CreatedAt.Local())
		j := i
		for j < len(entries) && entry.StartOfDay(entries[j].CreatedAt.Local()).Equal(day) {
			j++
		}
		fmt.Fprintln(w)
		header(w, day)
		writeEntries(w, entries[i:j])
		i = j
	}
	writeSummary(w, summary)
}

func Todos(w io.Writer, todos []entry.Todo) {
	if len(todos) == 0 {
		return
	}
	fmt.Fprintln(w)
	todoStyle.Fprintln(w, "Pendientes")
	for _, t := range todos {
		todoStyle.Fprintf(w, " [%d] %s  %s\n", t.ID, t.CreatedAt.Local().Format("02/01"), t.Text)
	}
}

func header(w io.Writer, day time.Time) {
	fmt.Fprintf(w, "%s, %s\n", weekdays[day.Weekday()], day.Format("02/01"))
}

func writeEntries(w io.Writer, entries []entry.Entry) {
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

func writeSummary(w io.Writer, summary string) {
	if summary == "" {
		return
	}
	fmt.Fprintln(w)
	summaryStyle.Fprintln(w, " Resumen: "+summary)
}
