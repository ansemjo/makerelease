// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

// output: release directory
var outdir string

// add flag to command
func addOutdirFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&outdir, "dir", "d", "release", "release directory")
}

func checkOutDirFlag(cmd *cobra.Command) (err error) {
	stat, err := os.Stat(outdir)
	if os.IsNotExist(err) || stat.IsDir() {
		return nil
	}
	if err != nil {
		return err
	}
	return errors.New("file exists: " + outdir)
}
