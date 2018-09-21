// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package main

import (
	"github.com/ansemjo/makerelease/cli/mkr"
	"github.com/ansemjo/makerelease/cli/tar"
	"github.com/spf13/cobra"
)

func init() {
	this := makeReleaseCmd

	// add to main, disable sorting
	mkrCmd.AddCommand(this)
	this.Flags().SortFlags = false

	// add flags
	addOutdirFlag(this)
	addInfileFlag(this)
	addTagFlag(this)
	addTargetsFlag(this)
	addEnvironmentFlag(this)

}

// make a release from passed source tarball
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
  git archive --prefix=./ HEAD | mkr rl -d output`,

	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return checkAll(cmd, checkTargetFlag, checkEnvironmentFlag, checkOutDirFlag, checkInFileFlag)
	},

	Run: func(cmd *cobra.Command, args []string) {

		// nest the function to always run deferred cleanup
		err := func() (err error) {

			// add current mkr version in environment
			environment = append(environment, "MKR_VERSION=mkr version "+mkrCmd.Version)

			// build the release
			cfg := mkr.MakeReleaseConfig{Targets: targets, Env: environment, Image: tag}
			release, err := mkr.MakeRelease(infile, cfg)
			if err != nil {
				return
			}
			defer release.Close()

			// untar to target directory, stripping the path prefix
			err = tar.Untar(outdir, release, "releases/")

			return
		}()

		// handle any error with nonzero exit code
		handleError(err)

	},
}
