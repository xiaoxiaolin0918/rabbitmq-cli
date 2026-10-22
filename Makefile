GO?=go
export PATH:=D:/tools/go/bin:$(PATH)
export GOPROXY?=goproxy.cn,direct

.PHONY: build test

build:
	$(GO) build -o rabbitmq.exe .

test:
	$(GO) test ./...
