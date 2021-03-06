# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

FROM golang:1.11.0-alpine

# install required and addon software
ARG SOFTWARE
RUN apk add --no-cache \
  bash \
  file \
  gcc \
  git \
  lzip \
  make \
  musl-dev \
  tar \
  upx \
  xz \
  ${SOFTWARE}

# prepare environment
ARG MKR_IMAGE=unknown
ENV \
  CACHEDIR=/tmp/buildcache \
  GOCACHE=/tmp/gocache \
  TEMPDIR=/tmp/makerelease \
  RELEASEDIR=/releases \
  WORKDIR=/build \
  MKR_IMAGE=${MKR_IMAGE} \
  MKR_VERSION=unknown

# add unpriviledged user
RUN adduser -u 1000 -D makerelease

# create workdir and release directory
RUN mkdir -m 1777 -p "${RELEASEDIR}" "${WORKDIR}"
WORKDIR ${WORKDIR}
USER makerelease

# entrypoint script, arguments can be given during docker run
COPY makerelease.sh /usr/bin/makerelease.sh
ENTRYPOINT [ "bash", "/usr/bin/makerelease.sh" ]
