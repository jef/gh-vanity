.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o sgv

.PHONY: clean
clean:
	rm -rf sgv

.PHONY: dist
dist:
	GOOS=windows GOARCH=amd64 go build -o stargazer-vanity.exe
	GOOS=linux GOARCH=amd64 go build -o stargazer-vanity.linux
	GOOS=darwin GOARCH=amd64 go build -o stargazer-vanity.darwin
