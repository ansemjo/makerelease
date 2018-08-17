IMAGE := go-releaser

RELEASES := $(PWD)/release
UIDGID := $(shell echo "$$(id -u):$$(id -g)")

.PHONY : default release image
default : image

# use targets after 'release' as arguments
# https://stackoverflow.com/a/14061796
ifeq (release,$(firstword $(MAKECMDGOALS)))
  RELEASE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RELEASE_ARGS):;@:)
endif

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
	docker build -t $(IMAGE) .

makerelease: $(shell ls cli/*.go)
	go generate cli/*.go
	packr
	CGO_ENABLED=0 go build -o $@ cli/*.go
	packr clean
	#upx $@
