FROM golang:1.10-alpine

# install required software
RUN apk add --no-cache \
  bash \
  file \
  git \
  lzip \
  make \
  upx \
  xz

# 
WORKDIR /build

# use unprivileged user
#RUN adduser -h /build -g 'Go Releaser' -D -k /dev/null builder
#USER builder

# prepare environment
ENV \
  GOCACHE=/tmp/go-build-cache \
  TEMPGOPATH=/tmp/go-build \
  RELEASES=/releases

# entrypoint script, arguments can be given during docker run
COPY release.sh /usr/bin/release.sh
ENTRYPOINT [ "bash", "/usr/bin/release.sh" ]