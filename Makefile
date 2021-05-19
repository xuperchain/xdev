
all: install

export GO111MODULE=on


build:
	ls bin || mkdir bin
	go build -o bin/xdev github.com/xuperchain/xdev

install:
	go install github.com/xuperchain/xdev

test:build
	go test ./...
	bin/xdev test jstest/testdata/jstest.test.js
lint:
	go vet ./...
