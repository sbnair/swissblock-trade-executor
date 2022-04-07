#!/bin/bash
MODULE=github.com/omhen/swissblock-trade-executor/v2
BINARY_NAME=trader

.PHONY: build
build: $(info Building project...)
	go build -v -a -o cli/$(BINARY_NAME) $(MODULE)/cli


.PHONY: test
test:
	go test -cover ./...

clean:
	rm ./cli/$(BINARY_NAME)

