# this makefile requires Go 1.11+
.PHONY : mkrelease-prepare mkrelease mkrelease-finish

mkrelease-prepare:
	mkdir -p assets
	tar cf assets/context.tar -C ../container/ .
	go get github.com/gobuffalo/packr/...
	go mod vendor

EXT := $(if $(findstring windows,$(OS)),.exe)
mkrelease:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
		packr build -o $(RELEASEDIR)/mkr-$(OS)-$(ARCH)$(EXT)

mkrelease-finish:
	cd $(RELEASEDIR) && sha256sum * | tee sha256sums
