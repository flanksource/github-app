
default: build
NAME:=github-app

ifeq ($(GITHUB_REF),)
GITHUB_REF := dev
endif
ifeq ($(VERSION),)
VERSION := $(shell echo $(GITHUB_REF) | sed 's|refs/tags/||')  built $(shell date)
endif

.PHONY: help
help:
	@cat docs/make-targets.md

.PHONY: release
release: setup linux darwin compress

.PHONY: setup
setup:
	which esc 2>&1 > /dev/null || go get -u github.com/mjibson/esc
	which github-release 2>&1 > /dev/null || go get github.com/aktau/github-release

.PHONY: build
build:
	go build -o ./.bin/$(NAME) -ldflags "-X \"main.version=$(VERSION)\""  main.go

.PHONY: linux
linux:
	GOOS=linux go build -o ./.bin/$(NAME) -ldflags "-X \"main.version=$(VERSION)\""  main.go

.PHONY: darwin
darwin:
	GOOS=darwin go build -o ./.bin/$(NAME)_osx -ldflags "-X \"main.version=$(VERSION)\""  main.go

.PHONY: compress
compress:
	upx -5 ./.bin/*

.PHONY: install
install:
	cp ./.bin/$(NAME) /usr/local/bin/

.PHONY: docker
docker:
	docker build ./ -t $(NAME)

#TODO: docs
#.PHONY: serve-docs
#serve-docs:
#	docker run --rm -it -p 8000:8000 -v $(PWD):/docs -w /docs squidfunk/mkdocs-material:5.1.5
#
#.PHONY: build-api-docs
#build-api-docs:
#	go run main.go docs api  pkg/types/config.go pkg/types/types.go pkg/types/nsx.go  > docs/reference/config.md
#	go run main.go docs cli "docs/cli"
#
#.PHONY: build-docs
#build-docs:
#	which mkdocs || pip install mkdocs mkdocs-material
#	mkdocs build -d build/docs
#
#.PHONY: deploy-docs
#deploy-docs:
#	which netlify 2>&1 > /dev/null || sudo npm install -g netlify-cli
#	netlify deploy --site b7d97db0-1bc2-4e8c-903d-6ebf3da18358 --prod --dir build/docs

.PHONY: lint
lint: build
	golangci-lint run --verbose --print-resources-usage
