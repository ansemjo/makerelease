// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

// miscellaneous utilities

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// simple check wether a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// run pre-run checks of cobra flags
func checkAll(cmd *cobra.Command, checker ...func(*cobra.Command) error) (err error) {
	for _, ch := range checker {
		err = ch(cmd)
		if err != nil {
			return
		}
	}
	return
}

// handle any non-nil errors
func handleError(err error) {
	if err != nil {

		// close open files
		if infile != nil {
			infile.Close()
		}

		// print error and exit
		fmt.Fprintln(os.Stderr, err)

		// suggest building the image
		if strings.HasPrefix(err.Error(), "Error: No such image") {
			fmt.Fprintln(os.Stderr, "try running 'mkr image' to build it")
		}

		os.Exit(1)

	}
}
