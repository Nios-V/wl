package cli

import (
	"os"
	"path/filepath"

	"github.com/Nios-V/wl/internal/config"
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
	root.AddCommand(newAddCmd(), newTodayCmd(), newYesterdayCmd(), newUndoCmd(),
		newTodoCmd(), newSummaryCmd(), newSprintCmd())
	return root
}

func wlDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".wl")
	return dir, os.MkdirAll(dir, 0o755)
}

func openStore() (*store.Store, error) {
	dir, err := wlDir()
	if err != nil {
		return nil, err
	}
	return store.Open(filepath.Join(dir, "wl.db"))
}

func loadConfig() (config.Config, error) {
	dir, err := wlDir()
	if err != nil {
		return config.Config{}, err
	}
	return config.Load(filepath.Join(dir, "config.yaml"))
}
