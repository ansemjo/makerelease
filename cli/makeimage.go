package main

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/spf13/cobra"
)

// simple command to generate the build image
var buildImageCmd = &cobra.Command{
	Use:     "image",
	Aliases: []string{"im"},
	Short:   "create the required docker image",
	Long:    "Build the required Docker image from the embedded context files locally.",
	Run: func(cmd *cobra.Command, args []string) {
		err := buildImage()
		handleError(err)
	},
}

func init() {
	cmd.AddCommand(buildImageCmd)
	buildImageCmd.Flags().SortFlags = false
	addTagFlag(buildImageCmd)
}

func buildImage() (err error) {

	// connect to docker daemon
	cli, ctx, err := newDockerClient()
	if err != nil {
		return
	}

	// open embedded build context
	buildcontext, err := assets.Open("context.tar")
	if err != nil {
		return
	}

	// begin building image
	build, err := cli.ImageBuild(ctx, buildcontext, types.ImageBuildOptions{
		Tags: []string{tag},
	})
	if err != nil {
		return
	}
	defer build.Body.Close()

	// print message stream to output and hadle errors
	err = jsonmessage.DisplayJSONMessagesStream(build.Body, os.Stdout, os.Stdout.Fd(), true, nil)
	if err != nil {
		if jerr, ok := err.(*jsonmessage.JSONError); ok {
			fmt.Fprintln(os.Stderr, jerr)
			os.Exit(jerr.Code)
		}
	}

	return
}
