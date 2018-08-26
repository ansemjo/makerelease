// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// a list of environment variables
var environment []string

// add flag to command
func addEnvironmentFlag(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&environment, "env", "e", []string{}, "environment variables to pass to container")
}

// sanity-check passed env variables
func checkEnvironmentFlag(cmd *cobra.Command) (err error) {
	for i, env := range environment {
		e := strings.Split(env, "=")
		// error on empty flag or env name
		if len(e) == 0 || e[0] == "" {
			return fmt.Errorf("invalid env flag: '%s'", env)
		}
		// if only name given, get from os.Env
		if len(e) == 1 {
			environment[i] = fmt.Sprintf("%s=%s", e[0], os.Getenv(e[0]))
		}
	}
	return
}
