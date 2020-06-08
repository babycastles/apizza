COVER=go tool cover

#VERSION=$(shell git describe --tags --abbrev=12)
VERSION=$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
GOFLAGS=-ldflags "-X $(shell go list)/cmd.version=$(VERSION)"

build: gen
	go build $(GOFLAGS)

install: gen
	go install $(GOFLAGS)

uninstall: clean
	go clean -i

test: test-build
	bash scripts/test.sh
	bash scripts/integration.sh ./bin/test-apizza
	@[ -d ./bin ] && [ -x ./bin/test-apizza ] && rm -rf ./bin

docker:
	docker build --rm -t apizza .

docker-test:
	docker build -f Dockerfile.test --rm -t apizza:$(VERSION) .
	docker run --rm -it apizza:$(VERSION)

release: gen
	scripts/release build

test-build: gen
	scripts/build.sh test

coverage.txt:
	@ echo '' > coverage.txt
	go test -v ./... -coverprofile=coverage.txt -covermode=atomic

html: coverage.txt
	$(COVER) -html=$<

gen:
	go generate ./...

clean:
	$(RM) -r coverage.txt release/apizza-* bin dist
	go clean -testcache
	go clean

all: test build release

.PHONY: install test clean html release gen docker
