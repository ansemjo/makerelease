# Copyright (c) 2018 Anton Semjonov
# Licensed under the MIT License

.PHONY : mkr install image dockerized release mkrelease-prepare mkrelease mkrelease-finish clean

VERSION := $(shell sed -n 's/^const version.*"\([0-9a-z.-]\+\)"$$/\1/p' cli/cmd_main.go)
IMAGE   := ansemjo/makerelease:$(VERSION)
RELEASE := $(PWD)/release

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
	docker build -t $(IMAGE) container/

# make a release in the container
MAKERID := mkr-$(shell date +%F-%s)
dockerized:
	docker run -i	--name $(MAKERID) $(IMAGE)
	docker cp $(MAKERID):/releases $(RELEASE)
	docker rm $(MAKERID)

# build the cli using the dockerized build process
release: image
	tar c ./makefile ./container ./cli | make dockerized

# delegate to submakefile
mkrelease-prepare mkrelease mkrelease-finish:
	make -C cli $@

# clean files not tracked by git
clean:
	git clean -fdx
