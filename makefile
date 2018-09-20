# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

VERSION := $(shell sed -n 's/^const version.*"\([0-9a-z.-]\+\)"$$/\1/p' cli/cmd_main.go)
IMAGE   := ansemjo/makerelease:$(VERSION)
RELEASE := $(PWD)/release

.PHONY : mkr install image dockerized

# compile binary, requires Go 1.11+
mkr: mkrelease-prepare mkrelease
	mv mkr-*-* mkr

# install in local path
PREFIX := $(shell [ $$(id -u) -eq 0 ] && echo /usr/local || echo ~/.local)
install: mkr
	install -d $(PREFIX)/bin
	install -m 755 $< $(PREFIX)/bin

# build docker image
image:
	docker build --build-arg MKR_IMAGE=$(IMAGE) -t $(IMAGE) container/

# make a release in the container
MAKERID := mkr-$(shell date --utc +%F-%s)
dockerized:
	docker run -i --name $(MAKERID) -e TARGETS=$(TARGETS) -e MKR_VERSION=makefile $(IMAGE)
	docker cp $(MAKERID):/releases $(RELEASE)
	docker rm $(MAKERID)

.PHONY: release selfrelease

# build the cli using the dockerized build process
release: clean image
	make -s archive | make dockerized

# release the cli using an intermediate mkr binary
selfrelease: clean image
	make -s archive | make dockerized TARGETS=host
	mv release/mkr-*-* mkr
	git clean -fdx -e mkr
	./mkr image
	make -s archive | ./mkr release

# delegate to submakefile
.PHONY: mkrelease-prepare mkrelease mkrelease-finish
mkrelease-prepare mkrelease mkrelease-finish:
	make -C cli $@

.PHONY: archive clean

# output a project archive of current HEAD
archive:
	@git archive --prefix=mkr/ HEAD

# clean files not tracked by git
clean:
	git clean -fdx
