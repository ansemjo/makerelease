// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// main cli command
var cmd = &cobra.Command{
	Use: "mkr",
}

func init() {
	cobra.EnableCommandSorting = false
}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func handleError(err error) {
	if err != nil {

		// close open files
		if infile != nil {
			infile.Close()
		}

		// print error and exit
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)

	}
}
