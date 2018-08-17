# this makefile requires Go 1.11+
.PHONY : prepare-release release finish-release

prepare-release:
	mkdir -p assets
	tar cf assets/context.tar -C ../container/ .
	go get github.com/gobuffalo/packr/...
	go mod vendor

release:
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) \
		packr build -o $(RELEASEDIR)/mkr-$(OS)-$(ARCH)

finish-release:
	cd $(RELEASEDIR) && sha256sum * | tee sha256sums
