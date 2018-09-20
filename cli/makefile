# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

# this makefile requires Go 1.11+
.PHONY : mkrelease-prepare mkrelease mkrelease-finish

# use host os/arch and current directory by default
OS    := $(shell go env GOHOSTOS)
ARCH  := $(shell go env GOHOSTARCH)
RELEASEDIR := $(PWD)

assets/context.tar.gz: $(shell find ../container -type f)
	mkdir -p assets
	(cd ../container && tar cz --mtime="1970-01-01 00:00:00 UTC" --owner=0 --group=0 *) > $@

mkrelease-prepare: assets/context.tar.gz
	go get github.com/gobuffalo/packr/...
	go mod download

EXT := $(if $(findstring windows,$(OS)),.exe)
mkrelease:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
		packr build -o $(RELEASEDIR)/mkr-$(OS)-$(ARCH)$(EXT)

mkrelease-finish:
	upx $(RELEASEDIR)/* || true
	printf "# built with %s in %s\n" "$$MKR_VERSION" "$$MKR_IMAGE" > $(RELEASEDIR)/SHA256SUMS
	cd $(RELEASEDIR) && sha256sum mkr-* | tee -a SHA256SUMS
