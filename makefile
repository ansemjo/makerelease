IMAGE := ansemjo/makerelease

RELEASES := $(PWD)/release
UIDGID := $(shell echo "$$(id -u):$$(id -g)")

.PHONY : default image dockerized-release prepare-release release finish-release
default : self

# create output directory
$(RELEASES):
	mkdir -p "$@"

# build docker image
image:
	docker build -t $(IMAGE) container/

# selfbuild binary releases
dockerized-release: $(RELEASES)
	@docker run \
		--init --rm -i \
		-v $(RELEASES):/releases \
		-u $(UIDGID) \
		$(IMAGE)

# build the cli using the dockerized process
self: image
	tar c ./makefile ./container ./cli | make dockerized-release

# delegate to submakefile
prepare-release release finish-release:
	make -C cli $@

# mkr: $(shell ls cli/*.go container/*)
# 	packr
# 	CGO_ENABLED=0 go build -o $@ cli/*.go
# 	packr clean
# 	upx $@
