package main

import (
	"os"

	"github.com/Nios-V/wl/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
