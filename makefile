IMAGE := ansemjo/makerelease

RELEASES := $(PWD)/release
UIDGID := $(shell echo "$$(id -u):$$(id -g)")

.PHONY : default releases image prepare-release release
default : image

# selfbuild binary releases
$(RELEASES):
	mkdir -p "$@"
releases: $(RELEASES)
	@docker run \
		--init --rm -i \
		--tmpfs /tmp:rw,nosuid,nodev,exec \
		-v $(RELEASES):/releases \
		-e TARGETS \
		-u $(UIDGID) \
		$(IMAGE) $(RELEASE_ARGS)

# targets for selfbuilding
prepare-release:
	mkdir -p cli/assets
	tar cf cli/assets/context.tar -C container/ .
	go get github.com/gobuffalo/packr/...
	go get golang.org/x/vgo
	packr
	cd cli && vgo mod vendor && ln -Ts vendor src

release-target-list:
	@echo linux/amd64 linux/arm64

release:
	cd cli && \
	CGO_ENABLED=0 GOPATH=$$PWD GOOS=$(OS) GOARCH=$(ARCH) \
	go build -o $(RELEASEDIR)/mkr-$(OS)-$(ARCH)

finish-release:
	cd $(RELEASEDIR) && sha256sum * | tee sha256sums

self: image $(RELEASES)
	tar cf sources.tar --exclude=sources.tar --exclude=.git .
	make releases < sources.tar

image:
	docker build -t $(IMAGE) container/

mkr: $(shell ls cli/*.go container/*)
	packr
	CGO_ENABLED=0 go build -o $@ cli/*.go
	packr clean
	upx $@
