package cli

import (
	"time"

	"github.com/Nios-V/wl/internal/entry"
	"github.com/Nios-V/wl/internal/render"
	"github.com/Nios-V/wl/internal/store"
	"github.com/spf13/cobra"
)

func newTodayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "today",
		Short: "Muestra las entradas de hoy",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showDay(cmd, time.Now())
		},
	}
}

func newYesterdayCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "yesterday",
		Short: "Muestra las entradas del día hábil anterior",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showDay(cmd, entry.PreviousWorkday(time.Now()))
		},
	}
}

func showDay(cmd *cobra.Command, day time.Time) error {
	from, to := entry.DayRange(day)
	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	entries, err := st.ListBetween(from, to)
	if err != nil {
		return err
	}
	summary, err := st.Summary(store.KindDay, from.Format(entry.DateLayout))
	if err != nil {
		return err
	}
	todos, err := st.OpenTodos()
	if err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	render.Day(out, from, entries, summary)
	render.Todos(out, todos)
	return nil
}
