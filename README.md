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

Currently, the container uses Go 1.11 RC2 and thus supports vendoring with `go mod`.

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

| variable     | available   | description                                                                                                           |
| ------------ | ----------- | --------------------------------------------------------------------------------------------------------------------- |
| `TIMESTAMP`  | all         | Starting time in `YYYY-MM-DD-UNIXEPOCH` format. _Careful:_ including<br> timestamps can often break reproduceability! |
| `WORKDIR`    | all         | Directory with extracted source code.                                                                                 |
| `TEMPDIR`    | all         | Temporary scratchpad.                                                                                                 |
| `RELEASEDIR` | all         | Output directory which is commonly mounted from the host machine.                                                     |
| `OS`         | `mkrelease` | Target operating system. Commonly `linux`, `darwin`, `*bsd`, `windows`, etc.                                          |
| `ARCH`       | `mkrelease` | Target architecture. Commonly `amd64`, `386`, `arm`, etc.                                                             |

### source code

The source code is piped into the container on stdin and needs to be in the format that many online
code repositories offer: there needs to be a single directory with the projects source code - or at
least the makefile - inside. During extraction, `tar` is called with `--strip-components=1` so that
source code lands in `WORKDIR` directly. For example, the downloaded archive for this project looks
like this:

```shell
$ tar tf makerelease-master.tar.gz
makerelease-master/
makerelease-master/.gitignore
makerelease-master/cli/
makerelease-master/cli/assets.go
[...]
```

If you want to pack a local directory be sure to use a single prefix component, even if it's just
`./`. I.e.:

```shell
~$ tar cf archive.tar -C /path/to/source/code ./
```

And **not** directly from within the directory:

```shell
/path/to/source/code$ tar cf ~/archive.tar cli/ container/ makefile [...]
```

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
instruct Docker to use the forwarded socket. Then continue with [mkr commands](#binary-usage).

### example makefile

For an example, take a look at this projects own [makefile](cli/makefile).

## building

As noted above, the Docker image is the actual heart of this project. You can use it completely
standalone by building and running it manually like this:

```shell
$ docker build -t makerelease container/
$ cat master.tar.gz | docker run --rm -i -v $PWD/release:/release makerelease
```

Or use the makefile target from within this directory:

```shell
$ make image
$ make dockerized-release < ~/master.tar.gz
```

### build the binary

Alternatively, compile and use the included Go binary `mkr` under [`cli/`](cli/). To compile and
install it use (`$GOPATH/bin` should be in your `PATH`):

```shell
$ make install
$ mkr --help
```

This currently requires that your `go` binary points to a Go release 1.11 or higher, which supports
the `mod` subcommand.

To build all the binary releases with the container (to build itself so to say):

```shell
$ make release
```

This will take a while and then place binaries in the `./release/` subdirectory.

## binary usage

The binary is more or less just an interface to start the [Docker](https://github.com/moby/moby)
container with the correct parameters and pipe the source code into it in a usable way, independent
from this project makefile.

If you have the binary on a machine where the Docker image does not yet exist, you can use the
binary to build it because the required build context is
[embedded](https://github.com/gobuffalo/packr) in the binary:

```shell
$ mkr image
```

Then, to build releases use the `release` / `rl` subcommand:

```shell
$ mkr rl -f master.tar.gz
```

Specify the output directory with `-d`:

```shell
$ cat master.tar.gz | mkr rl -d /path/to/release/output
```

Specify target list overrides with `-T`:

```shell
$ tar c ./ | mkr rl -T linux/amd64
```

If in doubt, you can use `--help` at any point. The CLI is built with the excellent
[cobra](https://github.com/spf13/cobra) commander which provides nice-looking usage information:

```shell
$ mkr help
$ mkr release --help
```
