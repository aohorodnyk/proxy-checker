 # Go parameters
GOCMD=go
BINARY_NAME=proxy-checker

build:
	$(GOCMD) build -o $(BINARY_NAME) src/*.go

run:
	$(GOCMD) run src/*.go

.PHONY: build
