SHELL:=/bin/bash

.DEFAULT_GOAL := default

name=tokens

build:
	@GOOS=darwin GOARCH=amd64 go build -o cmd/${name}-darwin-amd64
	@GOOS=darwin GOARCH=arm64 go build -o cmd/${name}-darwin-arm64
	@GOOS=linux GOARCH=arm64 go build -o cmd/${name}-linux-arm64
	@GOOS=linux GOARCH=amd64 go build -o cmd/${name}-linux-amd64
	@GOOS=linux GOARCH=386 go build -o cmd/${name}-linux-386
	@GOOS=windows GOARCH=amd64 go build -o cmd/${name}-windows-amd64.exe

default:
	@go build -o cmd/tokens
	@echo -e "Build done:\n$(pwd)cmd/${name}"

install:
	@chmod 777 cmd/tokens
	@echo -e '\nPATH+=":$(CURDIR)/cmd"\n' >> ${HOME}/.bashrc
