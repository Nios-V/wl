package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/Nios-V/wl/internal/entry"
	"github.com/Nios-V/wl/internal/store"
	"github.com/spf13/cobra"
)

func newSummaryCmd() *cobra.Command {
	var yesterday bool
	cmd := &cobra.Command{
		Use:   "summary <texto>",
		Short: "Guarda el resumen del día (reemplaza el anterior)",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			day := time.Now()
			if yesterday {
				day = entry.PreviousWorkday(day)
			}
			st, err := openStore()
			if err != nil {
				return err
			}
			defer st.Close()

			text := strings.Join(strings.Fields(strings.Join(args, " ")), " ")
			if err := st.SetSummary(store.KindDay, day.Format(entry.DateLayout), text, time.Now()); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Resumen del %s guardado.\n", day.Format("02/01"))
			return nil
		},
	}
	cmd.Flags().BoolVarP(&yesterday, "yesterday", "y", false, "guardarlo en el día hábil anterior")
	return cmd
}
