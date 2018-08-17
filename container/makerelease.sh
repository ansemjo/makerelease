#!/usr/bin/env bash

# fail on errors
set -e

# timestamp in YYYY-MM-DD-UNIXEPOCH format
export TIMESTAMP=$(date --utc +%F-%s)
printf 'starting makerelease.sh %s ...\n' "$TIMESTAMP"

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
DEFAULT_TARGETS=$(echo {linux,darwin,{free,net,open}bsd,plan9,windows}/{386,amd64} \
  {solaris,dragonfly}/amd64 {linux,darwin}/arm{,64} {free,net,open}bsd/arm )
MAKE_TARGETS=$(make -e mkrelease-targets) || true
TARGETS=${TARGETS:-${MAKE_TARGETS:-$DEFAULT_TARGETS}}
printf 'defined release targets:\n'; printf ' - %s\n' $TARGETS

# finally make targets
for target in $TARGETS; do

  printf 'make target: %s ...\n' "$target"

  OS=$(dirname "$target")
  ARCH=$(basename "$target")
  export OS ARCH

  make -e mkrelease

done

# finish up release, e.g. calculate checksums
printf 'finish up release ...\n'
make -e mkrelease-finish
