IMAGE := ansemjo/makerelease

RELEASES := $(PWD)/release
UIDGID := $(shell echo "$$(id -u):$$(id -g)")

.PHONY : default image install dockerized-release prepare-release release finish-release self clean
default : mkr

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
mkrelease-prepare mkrelease mkrelease-finish:
	make -C cli $@

# requires Go 1.11+
mkr: mkrelease-prepare
	cd cli && CGO_ENABLED=0 packr build -o ../$@

# install
PREFIX := $(shell [ $$(id -u) -eq 0 ] && echo /usr/local || echo ~/.local)
install: mkr
	install -m 755 $< $(PREFIX)/bin

# clean files not tracked by git
clean:
	git clean -fdx