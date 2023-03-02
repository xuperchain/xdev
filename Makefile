# set go env if needed
GO := go

all: build

export GO111MODULE=on

build:
	@ls bin 2>&1 >/dev/null || mkdir bin
	$(GO) build -o bin/xdev github.com/xuperchain/xdev

install:
	$(GO) install github.com/xuperchain/xdev

unit-test:
	$(GO) test ./...

build-test:build
	bin/xdev build -o testdata/counter-c.wasm testdata/counter.cc
	bin/xdev test testdata/counter.test.js

test:unit-test build-test

lint:
	$(GO) vet ./...

coverage:
	$(GO) test -coverprofile=coverage.txt -covermode=atomic ./...
