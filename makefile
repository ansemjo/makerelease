IMAGE := go-releaser
NAME := unknown

release:
	@docker run \
		--init --rm -i \
		--tmpfs /tmp:rw,nosuid,nodev,exec \
		-v $(PWD)/releases:/releases \
		-e TARGETS \
		$(IMAGE) $(NAME)

image:
	docker build -t $(IMAGE) .