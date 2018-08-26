# installing `mkr`

## from sources

To compile `mkr` from sources use:

    make mkr
    ./mkr --help

This currently requires that your `go` binary points to a Go release 1.11 or higher, which supports
the `mod` subcommand. It also uses [gobuffalo/packr](https://github.com/gobuffalo/packr) to embed
the container build context.

The manual steps would be:

    go get github.com/gobuffalo/packr/...
    (cd cli && CGO_ENABLED=0 packr build -o ../mkr)

You can install the resulting binary, assuming that `~/.local/bin` is in your `$PATH`:

    make install

To build all the binary releases with the container (to release itself so to say):

    make release

This will take a while and then place binaries in the `./release/` subdirectory.

## download a release

You can of course also download a release from GitHub. First, get my key from
[keybase](https://keybase.io/ansemjo):

    curl https://keybase.io/ansemjo/pgp_keys.asc | gpg --import

Download a binary, check the signature on the checksums and verify that the binary matches its
supposed hash:

    DOWNLOAD=https://github.com/ansemjo/makerelease/releases/download
    RELEASE=0.2.2 #<< substitue desired release
    curl -LO $DOWNLOAD/$RELEASE/mkr-linux-amd64
    curl -LO $DOWNLOAD/$RELEASE/SHA256SUMS.asc
    gpg --verify SHA256SUMS.asc
    sha256sum -c SHA256SUMS.asc --ignore-missing

Make it executable and see the [usage](README.md#usage) information:

    chmod +x mkr-linux-amd64
    ./mkr-linux-amd64 --help

Verify that you can build the exact same binaries yourself:

    curl -LO https://github.com/ansemjo/makerelease/archive/$RELEASE.tar.gz
    ./mkr-linux-amd64 rl -T linux/amd64 < $RELEASE.tar.gz
    sha256sum mkr-linux-amd64 release/mkr-linux-amd64

Now build other software compatible with the
[makefile requirements](README.md#your-project-must-bring-).
