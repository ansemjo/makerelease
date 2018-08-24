// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"github.com/ansemjo/makerelease/cli/docker"
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
'mkrelease-prepare', 'mkrelease' and 'mkrelease-finish'. During the build
the following environment variables can and should be used:

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
		err = checkTargetFlag(cmd)
		if err != nil {
			return
		}
		err = checkOutDirFlag(cmd)
		if err != nil {
			return
		}
		return checkInFileFlag(cmd)
	},

	Run: func(cmd *cobra.Command, args []string) {

		cfg := docker.MakeReleaseConfig{Targets: targets, Image: tag}
		release, err := docker.MakeRelease(infile, cfg)
		handleError(err)

		err = Untar(outdir, release, "releases/")
		if err != nil {
			return
		}

	},
}

func init() {
	cmd.AddCommand(makeReleaseCmd)
	makeReleaseCmd.Flags().SortFlags = false
	addOutdirFlag(makeReleaseCmd)
	addInfileFlag(makeReleaseCmd)
	addTagFlag(makeReleaseCmd)
	addTargetsFlag(makeReleaseCmd)
}
