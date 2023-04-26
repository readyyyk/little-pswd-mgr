.DEFAULT_GOAL := default

name=tokens

build:
	@GOOS=darwin GOARCH=386 go build -o cmd/{name}-darwin-386
	@GOOS=linux GOARCH=386 go build -o cmd/{name}-linux-386
	@GOOS=windows GOARCH=386 go build -o cmd/{name}-windows-386

default:
	@go build -o cmd/tokens
	@echo -e "Build done:\n$(pwd)cmd/${name}"

install:
	@chmod 777 cmd/tokens
	@echo -e '\nPATH+=":$(CURDIR)/cmd"\n' >> ~/.bashrc
