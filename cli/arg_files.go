// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// input: source tarball
var (
	infile        io.ReadCloser
	infileArg     string
	infileFlag    = []string{"src", "f", "", "source tarball (default stdin)"}
	addInfileFlag = func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&infileArg, infileFlag[0], infileFlag[1], infileFlag[2], infileFlag[3])
	}
)

// open the input file, use stdin if none given
func checkInFileFlag(cmd *cobra.Command) (err error) {
	if cmd.Flag(infileFlag[0]).Changed && infileArg != "-" {
		infile, err = os.Open(infileArg)
	} else {
		infile = os.Stdin
	}
	if err != nil {
		return
	}
	if terminal.IsTerminal(int(infile.(*os.File).Fd())) {
		return errors.New("refusing to read tar from tty")
	}
	return
}

// output: release directory
var (
	outdir        string
	outdirArg     string
	outdirFlag    = []string{"dir", "d", "release", "release directory"}
	addOutdirFlag = func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&outdirArg, outdirFlag[0], outdirFlag[1], outdirFlag[2], outdirFlag[3])
	}
)

// create output directory and transform to absolute path
func checkOutDirFlag(cmd *cobra.Command) (err error) {
	err = os.MkdirAll(outdirArg, 0755)
	if err != nil {
		return
	}
	outdir, err = filepath.Abs(outdirArg)
	return
}
