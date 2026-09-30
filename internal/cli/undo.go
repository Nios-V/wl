package cli

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"github.com/Nios-V/wl/internal/store"
	"github.com/spf13/cobra"
)

func newUndoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "undo",
		Short: "Borra la última entrada (pide confirmación)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			st, err := openStore()
			if err != nil {
				return err
			}
			defer st.Close()

			last, err := st.Last()
			if errors.Is(err, store.ErrNotFound) {
				fmt.Fprintln(out, "No hay entradas.")
				return nil
			}
			if err != nil {
				return err
			}

			fmt.Fprintf(out, "Borrar: %s  %s\n¿Confirmar? [s/N] ", last.CreatedAt.Local().Format("02/01 15:04"), last.Text)
			answer, _ := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
			switch strings.ToLower(strings.TrimSpace(answer)) {
			case "s", "si", "sí", "y", "yes":
			default:
				fmt.Fprintln(out, "Cancelado.")
				return nil
			}

			if err := st.Delete(last.ID); err != nil {
				return err
			}
			fmt.Fprintln(out, "Borrada.")
			return nil
		},
	}
}
