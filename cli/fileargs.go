// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// input: source tarball
var (
	infile        io.ReadCloser
	infileArg     string
	infileFlag    = []string{"src", "f", "source tarball"}
	addInfileFlag = func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&infileArg, infileFlag[0], infileFlag[1], "", infileFlag[2])
	}
)

// open the input file, use stdin if none given
func checkInFileFlag(cmd *cobra.Command) (err error) {
	if cmd.Flag(infileFlag[0]).Changed && infileArg != "-" {
		f, err := os.Open(infileArg)
		if err != nil {
			return err
		}
		infile = f
	} else {
		infile = os.Stdin
	}
	return
}

// output: release directory
var (
	outdir        string
	outdirArg     string
	outdirFlag    = []string{"dir", "d", "release directory"}
	addOutdirFlag = func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&outdirArg, outdirFlag[0], outdirFlag[1], "releases", outdirFlag[2])
	}
)

// check that the output directory exists and transform to absolute path
func checkOutDirFlag(cmd *cobra.Command) (err error) {
	stat, err := os.Stat(outdirArg)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("%s: not a directory", outdirArg)
	}
	absolute, err := filepath.Abs(outdirArg)
	if err != nil {
		return err
	}
	outdir = absolute
	return
}
