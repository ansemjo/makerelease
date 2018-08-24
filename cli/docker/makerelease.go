// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package docker

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"os/user"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

// MakeReleaseConfig is a struct holding information for the
// makerelease container. Give the name of the image to be used
// and possibly override the release targets.
type MakeReleaseConfig struct {
	Image   string
	Targets []string
}

// MakeRelease builds a release from a sourcecode tar with a makefile
// specifying the mkrelease* targets. The release is returned as a
// Reader of a tar file.
func MakeRelease(sourcecode io.Reader, config MakeReleaseConfig) (release io.ReadCloser, err error) {

	// connect to docker daemon
	cli, ctx, err := newDockerClient()
	if err != nil {
		return
	}

	// TODO: is this needed when we are copying the release?
	// get current user as UID:GID string
	id, err := func() (string, error) {
		cur, err := user.Current()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s:%s", cur.Uid, cur.Gid), nil
	}()
	if err != nil {
		return
	}

	// assemble container environment
	var env []string
	if len(config.Targets) > 0 {
		tg := fmt.Sprintf("TARGETS=%s", strings.Join(config.Targets, " "))
		env = append(env, tg)
	}

	// create the container
	c, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:     config.Image,
			OpenStdin: true,
			StdinOnce: true,
			User:      id,
			Env:       env,
		},
		&container.HostConfig{}, nil, "")
	if err != nil {
		return
	}

	// echo the container id
	fmt.Println("created container", c.ID[:12])

	// handle SIGINT / Ctrl-C / SIGKILL and remove container before exiting
	removeContainer := func() error { return cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{Force: true}) }
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, os.Kill)
	go func() {
		for range sigint {
			fmt.Fprintln(os.Stderr, "cancel")
			if err := removeContainer(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(1)
		}
	}()

	// start the container
	if err = cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		removeContainer()
		return
	}

	// attach to the container stdin/out/err
	attach, err := cli.ContainerAttach(ctx, c.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stderr: true,
		Stdout: true,
	})
	if err != nil {
		removeContainer()
		return
	}
	defer attach.Close()

	// copy input file to stdin of the container
	_, err = io.Copy(attach.Conn, sourcecode)
	if err != nil {
		removeContainer()
		return
	}
	attach.Conn.Close()

	// connect to the logging to follow progress
	logs, err := cli.ContainerLogs(ctx, c.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		removeContainer()
		return
	}

	// watch output
	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, logs)
	if err != nil {
		return
	}

	// inspect container state
	state, err := cli.ContainerInspect(ctx, c.ID)
	if err != nil {
		return
	}

	// and return any errors
	if state.State.ExitCode != 0 {
		err = errors.New(state.State.Error)
		return
	}

	// copy release files
	release, _, err = cli.CopyFromContainer(ctx, c.ID, "/releases/")
	if err != nil {
		return
	}

	err = removeContainer()
	if err != nil {
		return
	}

	err = cli.Close()
	return

}
