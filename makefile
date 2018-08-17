IMAGE := ansemjo/makerelease

RELEASES := $(PWD)/release
UIDGID := $(shell echo "$$(id -u):$$(id -g)")

.PHONY : default image install dockerized-release prepare-release release finish-release
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

# requires Go 1.11+
mkr: prepare-release
	cd cli && CGO_ENABLED=0 packr build -o ../$@

# install
PREFIX := ~/.local
install: mkr
	install -m 755 $< $(PREFIX)/bin