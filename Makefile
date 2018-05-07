ifeq ($(origin VERSION), undefined)
    VERSION := $(shell git describe --tags --always)
endif
# CHANEME
CONTAINER = arduima/admission-webhook
PKG = github.com/dkoshkin/admission-webhook

PORT = 443
HOST = 0.0.0.0

build:
	docker build -t $(CONTAINER) .
	docker tag $(CONTAINER) $(CONTAINER):$(VERSION)

test:
	@docker run                           \
		--rm                                \
		-e GLIDE_GOOS="linux"               \
		-u root:root                        \
		-v "$(shell pwd)":/go/src/$(PKG)    \
		-w /go/src/$(PKG)                   \
		arduima/golang-dev:1.10.2-alpine    \
		go test ./

.PHONY: vendor
vendor:
	@docker run                           \
		--rm                                \
		-v "$(shell pwd)":/go/src/$(PKG)    \
		-w /go/src/$(PKG)                   \
		arduima/glide:v0.13.1-1.10.2-alpine \
		install -v

.PHONY: update
update:
	@docker run                           \
		--rm                                \
		-v "$(shell pwd)":/go/src/$(PKG)    \
		-w /go/src/$(PKG)                   \
		arduima/glide:v0.13.1-1.10.2-alpine \
		up -v

push:
	docker push $(CONTAINER):$(VERSION)
	docker push $(CONTAINER):latest