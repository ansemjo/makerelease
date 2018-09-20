#!/usr/bin/env bash

# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

# fail on errors
set -e

# timestamp in YYYY-MM-DD-UNIXEPOCH format
export TIMESTAMP=$(date --utc +%F-%s)
printf 'starting makerelease.sh %s ...\n' "$TIMESTAMP"
printf 'running with %s in %s\n' \
  "${MKR_VERSION:-unknown mkr version}" \
  "${MKR_IMAGE:-unknown image}"

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
printf 'prepare release ...\n'
make -e mkrelease-prepare

# define target list in OS/ARCH format (env > make list > default)
DEFAULT_TARGETS=$(echo {linux,darwin,windows}/{386,amd64} linux/arm{,64} {free,open}bsd/{386,amd64,arm})
MAKE_TARGETS=$(make -se mkrelease-targets) || true
TARGETS=${TARGETS:-${MAKE_TARGETS:-$DEFAULT_TARGETS}}
printf 'defined release targets:\n'; printf ' - %s\n' $TARGETS

# finally make targets
for target in $TARGETS; do

  # build for host architecture
  if [[ $target == host ]]; then
    target=$(printf '%s/%s\n' $(go env GOHOSTOS GOHOSTARCH))
  fi

  printf 'make target: %s ...\n' "$target"

  OS=$(dirname "$target")
  ARCH=$(basename "$target")
  export OS ARCH

  make -e mkrelease

done

# finish up release, e.g. calculate checksums
printf 'finish up release ...\n'
make -e mkrelease-finish
