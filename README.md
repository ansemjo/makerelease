# makerelease

A utility to create reproducible builds by performing all build steps inside a Docker container.
Developed for Go but probably adaptable to other languages.

## the container provides ...

The heart of this project is the [makerelease.sh](/container/makerelease.sh) script and the
[dockerfile](/container/dockerfile) which puts it in a clean Golang container. Those two files
should be small enough to manually audit. The basic gist is:

[golang]: https://hub.docker.com/_/golang/

- `dockerfile`
  - create a new container from a [Golang] template with fixed version
  - install software commonly required for building
  - set up environment and paths like `WORKDIR`, `RELEASEDIR` and `TEMPDIR`
  - put the `makerelease.sh` script in place and mark it as the entrypoint
- `makerelease.sh`
  - read a tar archive with the project source code from stdin and extract
    - optionally perform decompression
  - perform preparational steps defined in the project's `makefile`
  - decide which targets to build for, optionally asking the `makefile` again
  - loop over the defined targets
    - set `OS` and `ARCH` environment variables
    - perform build steps defined in `makefile`
  - finish up release, e.g. calculate checksums

The container uses Go 1.11 and thus supports modules with `go mod`.

## your project must bring ...

### makefile

For the above steps to work, the project needs to provide a makefile which makes proper use of the
given environment variables and defines the correct targets. Those targets are:

| target              | required | description                                                                                                                                                                                                                                   |
| ------------------- | -------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `mkrelease-prepare` | ✔        | Perform any necessary preparations. For example,<br> fetch requirements and vendor them, compile common<br> static assets, create folders, etc.                                                                                               |
| `mkrelease-targets` | ✗        | Echo a list of targets to build in a `os/arch` format. The<br> output will be consumed directly, so make sure to silence<br> the command as much as possible. Example targets can<br> be `linux/amd64`, `darwin/amd64` or `openbsd/arm`, etc. |
| `mkrelease`         | ✔        | Build the binary for a single target and write the file to<br> the output directory in `RELEASEDIR`                                                                                                                                           |
| `mkrelease-finish`  | ✔        | Wrap up your release. Calculate checksums of your binaries,<br> compress them, include documentation, etc.                                                                                                                                    |

### environment

The most useful environment variables during the above procedure:

| variable      | available   | description                                                                                                           |
| ------------- | ----------- | --------------------------------------------------------------------------------------------------------------------- |
| `TIMESTAMP`   | all         | Starting time in `YYYY-MM-DD-UNIXEPOCH` format. _Careful:_ including<br> timestamps can often break reproduceability! |
| `WORKDIR`     | all         | Directory with extracted source code.                                                                                 |
| `TEMPDIR`     | all         | Temporary scratchpad.                                                                                                 |
| `RELEASEDIR`  | all         | Output directory which is commonly mounted from the host machine.                                                     |
| `OS`          | `mkrelease` | Target operating system. Commonly `linux`, `darwin`, `*bsd`, `windows`, etc.                                          |
| `ARCH`        | `mkrelease` | Target architecture. Commonly `amd64`, `386`, `arm`, etc.                                                             |
| `MKR_VERSION` | all         | Version of the `mkr` binary used.                                                                                     |
| `MKR_IMAGE`   | all         | Version / tag of the Docker image used.                                                                               |

### source code

The source code is piped into the container on stdin and needs to be in the format that many online
code repositories offer: there needs to be a single directory with the projects source code - or at
least the makefile - inside. During extraction, `tar` is called with `--strip-components=1` so that
source code lands in `WORKDIR` directly. For example, the downloaded archive for this project looks
like this:

    $ tar tf makerelease-master.tar.gz
    makerelease-master/
    makerelease-master/.gitignore
    makerelease-master/cli/
    makerelease-master/cli/assets.go
    [...]

If you want to pack a local directory be sure to use a single prefix component, even if it's just
`./`:

    tar cf archive.tar -C /path/to/source/code ./

Or use git's built-in `archive` command:

    git archive --prefix=./ HEAD > latest.tar

### docker daemon

You must be able to connect to a Docker daemon with default or current environment settings. That
means either adding your user to the `docker` group if you are not root (this is a
[security risk](https://docs.docker.com/install/linux/linux-postinstall/#manage-docker-as-a-non-root-user)
though) or configuring
[environment variables](https://docs.docker.com/engine/reference/commandline/cli/#environment-variables)
beforehand.

If you have a CoreOS host at your disposal you could do this, for example:

```shell
$ ssh -L ~/docker.sock:/var/run/docker.sock -Nf core@coreos.mydomain.tld
$ export DOCKER_HOST=unix://$HOME/docker.sock
```

This forwards the docker socket, puts ssh in the background and exports a `DOCKER_HOST` value to
instruct Docker to use the forwarded socket. Then continue with [mkr commands](#usage).

### example makefile

For an example, take a look at this projects own [makefile](cli/makefile).

## building a project

As noted above, the Docker image is the actual heart of this project and you can use it completely
standalone:

    docker build -t makerelease container/
    docker run -i --name myrelease makerelease < master.tar.gz
    docker cp myrelease:/releases release/
    docker rm myrelease

The [cli](#usage) is a lot easier to use though:

    mkr image
    mkr release < master.tar.gz

# installation

Installation notes can be found in [INSTALL.md](INSTALL.md).

# usage

The binary is more or less just an interface to start the [Docker](https://github.com/docker/docker)
container with the correct parameters and pipe the source code into it in a user-friendly way,
completely independent from this project's makefile.

If the Docker daemon does not have the correct image yet, you can use the binary to build it because
the required build context is [embedded](https://github.com/gobuffalo/packr):

    mkr image

Then, to build releases use the `release` / `rl` subcommand:

    mkr rl -f master.tar.gz

Specify the output directory with `-d`:

    mkr rl -d /path/to/release/output < master.tar.gz

Specify target list overrides with `-T`:

    tar c ./ | mkr rl -T linux/amd64

If in doubt, you can use `help`/`--help` at any point. The CLI is built with the excellent
[cobra](https://github.com/spf13/cobra) commander which provides nice-looking usage information:

    mkr help

## library usage

It should be possible to use `mkr` in your own Go applications. I'd be interested to hear how that
works out. :)

You'll need to import `github.com/ansemjo/makerelease/cli/mkr`. See the cli implementation for
examples:

- `mkr.BuildImage()` in [cmd_buildimage.go > buildImageCmd > Run](cli/cmd_buildimage.go#L34)
- `mkr.MakeRelease()` in [cmd_makerelease.go > makeReleaseCmd > Run](cli/cmd_makerelease.go#L63)

`MakeRelease()` returns a `ReadCloser` of the release tar archive. The Docker client connection is
only closed and the container removed when you close that reader, so don't forget to do that.
