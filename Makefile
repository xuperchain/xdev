build:
	go build github.com/xuperchain/xdev

install:
	go install github.com/xuperchain/xdev

test:install
	go test ./...
	xdev test jstest/testdata/jstest.test.js
lint:
	golint ./...
