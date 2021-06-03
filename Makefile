
all: build

export GO111MODULE=on

build:
	@ls bin 2>&1 >/dev/null || mkdir bin
	go build -o bin/xdev github.com/xuperchain/xdev

install:
	go install github.com/xuperchain/xdev

unit-test:
	go test ./...

build-test:build
	bin/xdev build -o testdata/counter-c.wasm testdata/counter.cc
	bin/xdev test testdata/counter.test.js

test:unit-test build-test

lint:
	go vet ./...

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
