#!/usr/bin/env bash

# fail on errors
set -e

# printf but on stderr
eprintf() { printf >&2 "$@"; }

# echo release name
name=${1:?release name required}
eprintf 'releasing %s ...\n' "$name"

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

# make preparations, don't fail if target does not exist
eprintf 'make preparations if necessary ...\n'
make prepare-release || [[ $? -eq 2 ]] && true

# define target list
#DEFAULT_TARGETS=$(echo {darwin,freebsd,linux,openbsd}/{386,amd64} linux/arm{,64})
DEFAULT_TARGETS=$(echo linux/{386,amd64})
TARGETS=${TARGETS:-$DEFAULT_TARGETS}

# make targets
for target in $TARGETS; do

  eprintf 'making target: %s ...\n' $target

  os=$(dirname $target)
  arch=$(basename $target)
  file="$name-$os-$arch"

  make release \
    GOOS=$os \
    GOARCH=$arch \
    TEMPDIR=$TEMPGOPATH \
    RELEASE=$RELEASES/$file

done

(cd $RELEASES && sha256sum * > sha256sums)