package main

//go:generate bash -c "mkdir -p assets && tar cvf assets/context.tar -C ../ dockerfile makerelease.sh"

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/gobuffalo/packr"
	"github.com/spf13/cobra"
)

var imagegen = &cobra.Command{
	Use: "build",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return buildImage()
	},
}

var box packr.Box

func init() {

	cmd.AddCommand(imagegen)

	box = packr.NewBox("assets")

}

func buildImage() (err error) {

	cli, ctx, err := newDockerClient()
	if err != nil {
		return
	}

	// open embedded build context
	dockercontext, err := box.Open("context.tar")
	if err != nil {
		return
	}

	build, err := cli.ImageBuild(ctx, dockercontext, types.ImageBuildOptions{
		Tags: []string{image},
	})
	if err != nil {
		return
	}
	defer build.Body.Close()

	// print message stream to output
	err = jsonmessage.DisplayJSONMessagesStream(build.Body, os.Stdout, os.Stdout.Fd(), true, nil)
	if err != nil {
		if jerr, ok := err.(*jsonmessage.JSONError); ok {
			fmt.Fprintln(os.Stderr, jerr)
			os.Exit(jerr.Code)
		}
	}

	return
}
