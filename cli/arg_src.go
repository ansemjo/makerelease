// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"errors"
	"io"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

// input: source tarball
var infile io.ReadCloser
var infileArg string

// add flag to command
func addInfileFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&infileArg, "src", "f", "", "source tarball (default stdin)")
}

// open the input file, use stdin if none given
func checkInFileFlag(cmd *cobra.Command) (err error) {
	if cmd.Flag("src").Changed && infileArg != "-" {
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
