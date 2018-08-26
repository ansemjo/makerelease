// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

// MakeReleaseConfig is a struct holding information for the
// makerelease container. Give the name of the image to be used
// and possibly override the release targets.
type MakeReleaseConfig struct {
	Image   string
	Env     []string
	Targets []string
}

// MakeRelease builds a release from a sourcecode tar with a makefile
// specifying the mkrelease* targets. The release is returned as a
// Reader of a tar file. The caller shall close the returned Reader.
func MakeRelease(sourcecode io.Reader, config MakeReleaseConfig) (release io.ReadCloser, err error) {

	// connect to docker daemon
	client, ctx, cancel, err := newDockerClient()
	if err != nil {
		return
	}

	// assemble container environment
	env := config.Env
	if len(config.Targets) > 0 {
		tg := fmt.Sprintf("TARGETS=%s", strings.Join(config.Targets, " "))
		env = append(env, tg)
	}

	// create the container
	ContainerConfig := &container.Config{
		Image:     config.Image,
		Env:       env,
		OpenStdin: true,
		StdinOnce: true,
	}
	maker, err := client.ContainerCreate(ctx, ContainerConfig, &container.HostConfig{}, nil, "")
	if err != nil {
		cancel()
		client.Close()
		return
	}

	// echo the container id
	id := maker.ID
	fmt.Println("created container", id[:12])

	// lambda function to remove created container and close client
	// connection in a new context with a short deadline
	cleanup := func() (err error) {
		timeout := 5 * time.Second
		deadline, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
		defer cancel()
		err = client.ContainerRemove(deadline, id, types.ContainerRemoveOptions{Force: true})
		client.Close()
		return
	}

	// defer cleanup, but only if we have an error. otherwise
	// do nothing and wait for cleanup upon closing the CopyCloser.
	defer func() {
		if err != nil {
			clerr := cleanup()
			if clerr != nil {
				fmt.Fprintln(os.Stderr, clerr)
			}
		}
	}()

	// handle SIGINT / Ctrl-C to cancel running operations
	// TODO: lingering goroutine?
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		for range sigint {
			fmt.Fprintln(os.Stderr, "ancel")
			cancel()
		}
	}()

	// start the container
	err = client.ContainerStart(ctx, id, types.ContainerStartOptions{})
	if err != nil {
		return
	}

	// attach to the container stdin/out/err
	attach, err := client.ContainerAttach(ctx, id, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stderr: true,
		Stdout: true,
	})
	if err != nil {
		return
	}
	defer attach.Close()

	// copy input file to stdin of the container in the background
	go func() {
		_, err = io.Copy(attach.Conn, sourcecode)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			cancel()
			return
		}
		attach.Conn.Close()
	}()

	// connect to the logging to follow progress
	logs, err := client.ContainerLogs(ctx, id, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		return
	}
	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, logs)
	if err != nil {
		return
	}

	// inspect container state after exit and return any errors
	state, err := client.ContainerInspect(ctx, id)
	if err != nil {
		return
	} else if state.State.ExitCode != 0 {
		err = errors.New(state.State.Error)
		return
	}

	// copy release tarball
	releaseReader, _, err := client.CopyFromContainer(ctx, id, "/releases")
	if err != nil {
		return
	}

	// return a ~ReadCloser struct with cleanup functions
	release = CopyCloser{releaseReader, cleanup}

	return
}

// CopyCloser is a Reader returned from CopyFromContainer with
// some cleanup procedures to close the connection etc. upon closing
// the ReadCloser.
type CopyCloser struct {
	io.ReadCloser
	cleanup func() error
}

// Close closes the ReadCloser and then runs further cleanup steps.
func (r CopyCloser) Close() (err error) {
	r.ReadCloser.Close()
	return r.cleanup()
}
