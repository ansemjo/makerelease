// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"github.com/ansemjo/makerelease/cli/assets"
	"github.com/ansemjo/makerelease/cli/mkr"
	"github.com/spf13/cobra"
)

func init() {
	this := buildImageCmd

	// add to main, disable sorting
	cmd.AddCommand(this)
	this.Flags().SortFlags = false

	// add flags
	addTagFlag(this)

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
		err = mkr.BuildImage(bc, tag)
		handleError(err)

	},
}
