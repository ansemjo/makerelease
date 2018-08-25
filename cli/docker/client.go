// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package docker

import (
	"context"

	Docker "github.com/docker/docker/client"
)

// create default client and negotiate api version
func newDockerClient() (client *Docker.Client, ctx context.Context, cancel context.CancelFunc, err error) {
	ctx, cancel = context.WithCancel(context.Background())
	client, err = Docker.NewClientWithOpts(Docker.FromEnv)
	if err != nil {
		return
	}
	client.NegotiateAPIVersion(ctx)
	return
}
