
all: build

export GO111MODULE=on

build:
	@ls bin 2>&1 >/dev/null || mkdir bin
	go build -o bin/xdev github.com/xuperchain/xdev

install:
	go install github.com/xuperchain/xdev

test:build
	# go test 
	go test ./...
	bin/xdev build -o testdata/counter-c.wasm testdata/counter.cc
	bin/xdev test testdata/counter.test.js
lint:
	go vet ./...

coverage:
	go test -coverprofile=coverage.txt -covermode=atomic ./...
