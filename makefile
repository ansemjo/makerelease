IMAGE := go-releaser
NAME := unknown

.PHONY : default release image
default : image

# use targets after 'release' as arguments
# https://stackoverflow.com/a/14061796
ifeq (release,$(firstword $(MAKECMDGOALS)))
  RELEASE_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RELEASE_ARGS):;@:)
endif

release:
	@docker run \
		--init --rm -i \
		--tmpfs /tmp:rw,nosuid,nodev,exec \
		-v $(PWD)/releases:/releases \
		-e TARGETS \
		$(IMAGE) $(RELEASE_ARGS)

image:
	docker build -t $(IMAGE) .