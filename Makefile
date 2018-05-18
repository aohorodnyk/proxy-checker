 # Go parameters
GOCMD=go

ifeq ($(OS),Windows_NT)
	BINARY_NAME=proxy-checker.exe
else
	BINARY_NAME=proxy-checker
endif

build:
	$(GOCMD) build -o $(BINARY_NAME) main.go

run:
	$(GOCMD) run main.go

.PHONY: build
