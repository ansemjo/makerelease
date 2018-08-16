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
  CACHEDIR=/tmp/buildcache \
  GOCACHE=/tmp/gocache \
  TEMPDIR=/tmp/makerelease \
  RELEASEDIR=/releases \
  WORKDIR=/build \
  HOME=/build

# create workdir and release directory
RUN mkdir -p "${RELEASEDIR}" "${WORKDIR}" \
  && chmod 1777 "${RELEASEDIR}" "${WORKDIR}"
WORKDIR ${WORKDIR}

# entrypoint script, arguments can be given during docker run
COPY makerelease.sh /usr/bin/makerelease.sh
ENTRYPOINT [ "bash", "/usr/bin/makerelease.sh" ]
