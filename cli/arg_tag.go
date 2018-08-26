// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"github.com/spf13/cobra"
)

// the tagged docker image to use
var tag string

// add flag to command
func addTagFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&tag, "tag", "t", "ansemjo/makerelease:"+version, "docker image/tag to use")
}
