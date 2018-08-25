// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"github.com/ansemjo/makerelease/cli/assets"
	"github.com/ansemjo/makerelease/cli/docker"
	"github.com/spf13/cobra"
)

func init() {
	cmd.AddCommand(buildImageCmd)
	buildImageCmd.Flags().SortFlags = false
	addTagFlag(buildImageCmd)
}

// simple command to generate the required builder image
var buildImageCmd = &cobra.Command{

	Use:     "image",
	Aliases: []string{"im"},
	Short:   "create the required docker image",
	Long:    "Build the required Docker image from embedded context files.",

	Run: func(cmd *cobra.Command, args []string) {

		// open embedded build context
		bc, err := assets.Box.Open("context.tar")
		handleError(err)
		defer bc.Close()

		// build the image
		err = docker.BuildImage(bc, tag)
		handleError(err)

	},
}
