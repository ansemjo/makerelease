// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"os"

	"github.com/spf13/cobra"
)

// main cli command
var cmd = &cobra.Command{
	Use:     "mkr",
	Version: "0.1.1",
}

func init() {
	cobra.EnableCommandSorting = false
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
