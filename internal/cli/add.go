package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Nios-V/wl/internal/entry"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <texto>",
		Short: "Agrega una entrada",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return addEntry(cmd, args)
		},
	}
}

func addEntry(cmd *cobra.Command, args []string) error {
	text, tags := entry.ParseTags(strings.Join(args, "  "))
	if text == "" {
		return errors.New("la entrada no tiene texto")
	}
	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	id, err := st.Add(entry.Entry{CreatedAt: time.Now(), Text: text, Tags: tags})
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Guardada (id %d)\n", id)
	return nil
}
