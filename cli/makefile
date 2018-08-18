# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

# this makefile requires Go 1.11+
.PHONY : mkrelease-prepare mkrelease mkrelease-finish

mkrelease-prepare:
	mkdir -p assets
	(cd ../container && tar c --mtime=1970-01-01 *) > assets/context.tar
	go get github.com/gobuffalo/packr/...
	go mod vendor

EXT := $(if $(findstring windows,$(OS)),.exe)
mkrelease:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
		packr build -o $(RELEASEDIR)/mkr-$(OS)-$(ARCH)$(EXT)

mkrelease-finish:
	upx $(RELEASEDIR)/* || true
	cd $(RELEASEDIR) && sha256sum * | tee sha256sums
