#!/usr/bin/env bash

# fail on errors
set -e

# do not fail on missing make target
ignore-missing-target() { [[ $? -eq 2 ]] && true; }

# printf but on stderr
eprintf() { printf >&2 "$@"; }

# timestamp in YYYY-MM-DD-UNIXEPOCH format
export TIMESTAMP=$(date --utc +%F-%s)
export RELEASEDIR="$RELEASEDIR/$TIMESTAMP"
mkdir -p "$RELEASEDIR"
eprintf 'releasing @ %s ...\n' "$TIMESTAMP"

# unpack source tarball, with decompression based on mime-type
eprintf 'reading tar archive from stdin ...\n'
SOURCE=/tmp/source.tar
cat > $SOURCE
mime=$(file $SOURCE --mime-type | awk '{print $2}')
cat $SOURCE | (
  case $mime in
    'application/x-gzip') gzip -d ;;
    'application/x-bzip2') bzip2 -d ;;
    'application/x-lzip') lzip -d ;;
    'application/x-xz') xz -d ;;
    *) cat ;;
  esac
) | tar x --strip-components=1

# make any required preparations
eprintf 'make preparations if necessary ...\n'
make -e prepare-release || ignore-missing-target

# define target list in OS/ARCH format (env > make list > default)
DEFAULT_TARGETS=$(echo {darwin,freebsd,linux,openbsd}/{386,amd64} linux/arm{,64})
MAKE_TARGETS=$(make -e release-target-list) || ignore-missing-target
TARGETS=${TARGETS:-${MAKE_TARGETS:-$DEFAULT_TARGETS}}
eprintf 'defined release targets:\n'; eprintf ' - %s\n' $TARGETS

# finally make targets
for target in $TARGETS; do

  eprintf 'make target: %s ...\n' "$target"

  OS=$(dirname "$target")
  ARCH=$(basename "$target")
  export OS ARCH

  make -e release

done

# finish up release, e.g. calculate checksums
eprintf 'finish up release ...\n'
make -e finish-release
