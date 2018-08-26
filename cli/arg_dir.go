// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// output: release directory
var outdir string
var outdirArg string

// add flag to command
func addOutdirFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outdirArg, "dir", "d", "release", "release directory")
}

// create output directory and transform to absolute path
func checkOutDirFlag(cmd *cobra.Command) (err error) {
	err = os.MkdirAll(outdirArg, 0755)
	if err != nil {
		return
	}
	outdir, err = filepath.Abs(outdirArg)
	return
}
