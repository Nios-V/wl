package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Nios-V/wl/internal/entry"
	"github.com/Nios-V/wl/internal/render"
	"github.com/Nios-V/wl/internal/store"
	"github.com/spf13/cobra"
)

func newSprintCmd() *cobra.Command {
	var prev bool

	// selected devuelve el rango del sprint actual, o del anterior si se pasó --prev.
	selected := func() (time.Time, time.Time, error) {
		start, end, err := sprintRange(time.Now())
		if err != nil || !prev {
			return start, end, err
		}
		return sprintRange(start.AddDate(0, 0, -1))
	}

	cmd := &cobra.Command{
		Use:   "sprint",
		Short: "Muestra los logs del sprint y su resumen",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			start, end, err := selected()
			if err != nil {
				return err
			}
			st, err := openStore()
			if err != nil {
				return err
			}
			defer st.Close()

			entries, err := st.ListBetween(start, end)
			if err != nil {
				return err
			}
			summary, err := st.Summary(store.KindSprint, start.Format(entry.DateLayout))
			if err != nil {
				return err
			}
			render.Sprint(cmd.OutOrStdout(), start, end, entries, summary)
			return nil
		},
	}
	cmd.PersistentFlags().BoolVarP(&prev, "prev", "p", false, "usar el sprint anterior")

	cmd.AddCommand(&cobra.Command{
		Use:   "summary <texto>",
		Short: "Guarda el resumen del sprint (reemplaza el anterior)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			start, _, err := selected()
			if err != nil {
				return err
			}
			st, err := openStore()
			if err != nil {
				return err
			}
			defer st.Close()

			text := strings.Join(strings.Fields(strings.Join(args, " ")), " ")
			if err := st.SetSummary(store.KindSprint, start.Format(entry.DateLayout), text, time.Now()); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Resumen del sprint del %s guardado.\n", start.Format("02/01"))
			return nil
		},
	})
	return cmd
}

func sprintRange(t time.Time) (time.Time, time.Time, error) {
	cfg, err := loadConfig()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	if cfg.Sprint.Start == "" {
		return time.Time{}, time.Time{}, errors.New("falta sprint.start en ~/.wl/config.yaml (ej: start: 2026-09-28)")
	}
	anchor, err := time.ParseInLocation(entry.DateLayout, cfg.Sprint.Start, time.Local)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("sprint.start inválido: %w", err)
	}
	start, end := entry.SprintRange(anchor, t, cfg.Sprint.Days)
	return start, end, nil
}
