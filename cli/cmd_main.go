// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

const version = "0.3.4"

// main cli command
var mkrCmd = &cobra.Command{
	Use:     "mkr",
	Version: fmt.Sprintf("%s (%s/%s, runtime %s)", version, runtime.GOOS, runtime.GOARCH, runtime.Version()),
}

func init() {
	cobra.EnableCommandSorting = false
}

func main() {
	if err := mkrCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
