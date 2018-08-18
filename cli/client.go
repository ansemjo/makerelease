package main

import (
	"context"

	"github.com/docker/docker/client"
)

// create default client and negotiate api version
func newDockerClient() (cli *client.Client, ctx context.Context, err error) {
	ctx = context.Background()
	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return
	}
	cli.NegotiateAPIVersion(ctx)
	return
}
