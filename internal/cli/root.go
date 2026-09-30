package cli

import (
	"os"
	"path/filepath"

	"github.com/Nios-V/wl/internal/store"
	"github.com/spf13/cobra"
)

func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "wl",
		Short:        "Bitácora de trabajo en la terminal",
		Args:         cobra.ArbitraryArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return addEntry(cmd, args)
		},
	}
	root.CompletionOptions.DisableDefaultCmd = true
	root.AddCommand(newAddCmd(), newTodayCmd(), newYesterdayCmd(), newUndoCmd())
	return root
}

func openStore() (*store.Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".wl")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return store.Open(filepath.Join(dir, "wl.db"))
}
