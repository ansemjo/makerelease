FROM golang:1.10-alpine

# install required and addon software
ARG SOFTWARE
RUN apk add --no-cache \
  bash \
  file \
  git \
  lzip \
  make \
  upx \
  xz \
  ${SOFTWARE}

# prepare environment
ENV \
  GOCACHE=/tmp/go-build-cache \
  TEMPDIR=/tmp/makerelease \
  RELEASES=/releases \
  WORKDIR=/build

# create workdir and release directory
RUN mkdir -p "${RELEASES}" "${WORKDIR}" \
  && chmod 1777 "${RELEASES}" "${WORKDIR}"
WORKDIR /build

# entrypoint script, arguments can be given during docker run
COPY makerelease.sh /usr/bin/makerelease.sh
ENTRYPOINT [ "bash", "/usr/bin/makerelease.sh" ]