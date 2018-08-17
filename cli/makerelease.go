package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"os/user"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/cobra"
)

// main cli command
var makeReleaseCmd = &cobra.Command{
	Use:     "release",
	Aliases: []string{"rl"},
	Short:   "build a release from source code",
	Long: `Build a release from source code inside a Docker container.

A tar archive shall be given as input, either via stdin or with the '-f'
flag. It can be compressed with gzip, bzip2, lzip or xz. The first component
will be stripped during extraction and the resulting directory will be the
root for all build operations. Archives downloaded from online repositories
like GitHub or GitLab often conform to this format.

In this build root a Makefile is expected with at least the targets
'prepare-release', 'release' and 'finish-release'. During the build the
following environment variables can and should be used:

  RELEASEDIR  - the output directory for finished files
  WORKDIR     - the build root with extracted sources
  OS          - the target operating system (linux, freebsd, darwin, ...)
  ARCH        - the target architecture (amd64, 386, arm, ...)

The latter two are only available for the 'release' target.`,
	Example: `
Create the Docker image locally first:
  mkr image

Build from a downloaded source archive:
  mkr rl -f master.tar.gz

Pack a local code directory and pipe it directly:
  tar c -C /path/to/code ./ | mkr rl -d output`,
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		err = checkOutDirFlag(cmd)
		if err != nil {
			return
		}
		return checkInFileFlag(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := makeRelease(infile, outdir)
		handleError(err)
	},
}

func init() {
	cmd.AddCommand(makeReleaseCmd)
	makeReleaseCmd.Flags().SortFlags = false
	addOutdirFlag(makeReleaseCmd)
	addInfileFlag(makeReleaseCmd)
	addTagFlag(makeReleaseCmd)
}

// make a release from the sourcecode in the source tarball
func makeRelease(tar io.ReadCloser, releases string) (err error) {

	// connect to docker daemon
	cli, ctx, err := newDockerClient()
	if err != nil {
		return
	}

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

	// create the container
	c, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image:     tag,
			OpenStdin: true,
			StdinOnce: true,
			User:      id,
		},
		&container.HostConfig{
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
		return
	}

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
		return err
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
		return err
	}
	defer attach.Close()

	// copy input file to stdin of the container
	if _, err = io.Copy(attach.Conn, tar); err != nil {
		removeContainer()
		return err
	}

	// input is done, close
	tar.Close()
	attach.Conn.Close()

	// connect to the logging to follow progress
	logs, err := cli.ContainerLogs(ctx, c.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     true,
	})
	if err != nil {
		removeContainer()
		return err
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
		return errors.New(state.State.Error)
	}

	return cli.Close()

}
