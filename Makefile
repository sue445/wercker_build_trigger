# Requirements: git, go, go dep
NAME     := wercker_build_trigger
VERSION  := $(shell cat VERSION)
REVISION := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

.DEFAULT_GOAL := bin/$(NAME)

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: clean
clean:
	rm -rf bin/*
	rm -rf dist/*
	rm -rf vendor/

vendor:
	dep ensure

.PHONY: package
package: vendor
	for os in darwin linux windows; do \
		if [ $$os = "windows" ]; then \
			exefile="$(NAME).exe" ; \
		else \
			exefile="$(NAME)" ; \
		fi ; \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o dist/$${os}_$${arch}/$${exefile} ; \
			cd dist/$${os}_$${arch} ; \
			zip ../$(NAME)_$(VERSION)_$${os}_$${arch}.zip $${exefile} ; \
			cd ../.. ; \
		done; \
	done

.PHONY: install
install:
	go install $(LDFLAGS)

.PHONY: test
test:
	go test

.PHONY: tag
tag:
	git tag -a $(VERSION) -m "Release v$(VERSION)"
	git push --tags
