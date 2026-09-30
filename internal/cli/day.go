package cli

import (
	"time"

	"github.com/Nios-V/wl/internal/entry"
	"github.com/Nios-V/wl/internal/render"
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
	render.Day(cmd.OutOrStdout(), from, entries)
	return nil
}
