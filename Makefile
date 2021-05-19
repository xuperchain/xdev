
all: install

export GO111MODULE=on
export PATH = ${GOPATH}/bin:$(PATH)


build:
	go build github.com/xuperchain/xdev

install:
	go install github.com/xuperchain/xdev

test:install
	go test ./...
	xdev test jstest/testdata/jstest.test.js
lint:
	go vet ./...
