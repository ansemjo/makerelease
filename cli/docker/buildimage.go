// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package docker

import (
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/jsonmessage"
)

// BuildImage builds a Docker image from the passed build context.
func BuildImage(buildcontext io.Reader, image string) (err error) {

	// connect to docker daemon
	cli, ctx, err := newDockerClient()
	if err != nil {
		return
	}

	// begin building image
	build, err := cli.ImageBuild(ctx, buildcontext, types.ImageBuildOptions{Tags: []string{image}})
	if err != nil {
		return
	}
	defer build.Body.Close()

	// print message stream to stdout
	err = jsonmessage.DisplayJSONMessagesStream(build.Body, os.Stdout, os.Stdout.Fd(), true, nil)
	return

}
