#!/usr/bin/env bash

# fail on errors
set -e

# do not fail on missing make target
ignore-missing-target() { [[ $? -eq 2 ]] && true || false; }

# timestamp in YYYY-MM-DD-UNIXEPOCH format
export TIMESTAMP=$(date --utc +%F-%s)
export RELEASEDIR="$RELEASEDIR/$TIMESTAMP"
mkdir -p "$RELEASEDIR"
printf 'releasing @ %s ...\n' "$TIMESTAMP"

# unpack source tarball, with decompression based on mime-type
printf 'reading tar archive from stdin ...\n'
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
printf 'make preparations if necessary ...\n'
make -e prepare-release || ignore-missing-target

# define target list in OS/ARCH format (env > make list > default)
DEFAULT_TARGETS=$(echo {darwin,freebsd,linux,openbsd}/{386,amd64} linux/arm{,64})
MAKE_TARGETS=$(make -e release-target-list) || ignore-missing-target
TARGETS=${TARGETS:-${MAKE_TARGETS:-$DEFAULT_TARGETS}}
printf 'defined release targets:\n'; printf ' - %s\n' $TARGETS

# finally make targets
for target in $TARGETS; do

  printf 'make target: %s ...\n' "$target"

  OS=$(dirname "$target")
  ARCH=$(basename "$target")
  export OS ARCH

  make -e release

done

# finish up release, e.g. calculate checksums
printf 'finish up release ...\n'
make -e finish-release || ignore-missing-target
