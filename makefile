IMAGE := ansemjo/makerelease

RELEASES := $(PWD)/release
UIDGID := $(shell echo "$$(id -u):$$(id -g)")

.PHONY : default release image
default : image

release: $(RELEASES)
	@docker run \
		--init --rm -i \
		--tmpfs /tmp:rw,nosuid,nodev,exec \
		-v $(RELEASES):/releases \
		-e TARGETS \
		-u $(UIDGID) \
		$(IMAGE) $(RELEASE_ARGS)

$(RELEASES):
	mkdir -p "$@"

image:
	docker build -t $(IMAGE) container/

mkr: $(shell ls cli/*.go container/*)
	tar cf cli/assets/context.tar -C container/ .
	packr
	CGO_ENABLED=0 go build -o $@ cli/*.go
	packr clean
	upx $@
