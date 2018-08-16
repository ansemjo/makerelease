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

var (
	infileFlag = []string{"src", "f", "source tarball"}
	infileArg  string
	infile     io.ReadCloser

	outdirFlag = []string{"dir", "d", "releases directory"}
	outdirArg  string
)

func addFileArgFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&infileArg, infileFlag[0], infileFlag[1], "", infileFlag[2])
	cmd.Flags().StringVarP(&outdirArg, outdirFlag[0], outdirFlag[1], "releases", outdirFlag[2])
}

func checkFileArgFlags(cmd *cobra.Command) (err error) {
	err = checkInFileFlag(cmd)
	if err != nil {
		return
	}
	err = checkOutDirFlag(cmd)
	return
}

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
	outdirArg = absolute
	return
}
