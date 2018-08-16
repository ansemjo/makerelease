package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// create default client and negotiate api version
func newDockerClient() (cli *client.Client, ctx context.Context, err error) {
	ctx = context.Background()
	cli, err = client.NewClientWithOpts()
	if err != nil {
		return
	}
	cli.NegotiateAPIVersion(ctx)
	fmt.Println("negotiated api version", cli.ClientVersion())
	return
}

func makeRelease(tar io.ReadCloser, releases string) {

	cli, ctx, err := newDockerClient()

	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image:     "go-releaser",
		OpenStdin: true,
		StdinOnce: true,
		User:      "1000:100",
	}, &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: releases,
				Target: "/releases",
			},
		},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	// handle SIGINT / Ctrl-C and kill container
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	go func() {
		for range sigint {
			cli.ContainerRemove(ctx, res.ID, types.ContainerRemoveOptions{Force: true})
		}
	}()

	if err = cli.ContainerStart(ctx, res.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	hj, err := cli.ContainerAttach(ctx, res.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stderr: true,
		Stdout: true,
	})
	if err != nil {
		panic(err)
	}
	defer hj.Close()

	if _, err = io.Copy(hj.Conn, tar); err != nil {
		panic(err)
	}
	tar.Close()
	hj.Conn.Close()

	out, err := cli.ContainerLogs(ctx, res.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

}

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
