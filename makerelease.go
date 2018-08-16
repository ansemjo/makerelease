package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func main() {

	ctx := context.Background()

	cli, err := client.NewClientWithOpts()
	if err != nil {
		panic(err)
	}

	cli.NegotiateAPIVersion(ctx)
	fmt.Println("negotiated api version", cli.ClientVersion())

	aenker, _ := os.Open("./aenker.tar")

	cli.ContainerRemove(ctx, "rel", types.ContainerRemoveOptions{
		Force: true,
	})

	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image:     "go-releaser",
		OpenStdin: true,
		StdinOnce: true,
	}, &container.HostConfig{
		AutoRemove: true,
	}, nil, "rel")
	if err != nil {
		panic(err)
	}

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

	if _, err = io.Copy(hj.Conn, aenker); err != nil {
		panic(err)
	}
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
