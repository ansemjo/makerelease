package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

var cmd = &cobra.Command{
	Use:   "makerelease",
	Short: "make reproducible releases",
	Run: func(cmd *cobra.Command, args []string) {
		demorelease()
	},
}

func init() {

}

func demorelease() {

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

	releases, _ := filepath.Abs("releases")

	res, err := cli.ContainerCreate(ctx, &container.Config{
		Image:     "go-releaser",
		OpenStdin: true,
		StdinOnce: true,
	}, &container.HostConfig{
		AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: releases,
				Target: "/releases",
			},
		},
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

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
