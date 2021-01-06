.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o sgv

.PHONY: clean
clean:
	rm -rf sgv

.PHONY: dist
dist:
	GOOS=windows GOARCH=amd64 go build -o svg.exe
	GOOS=linux GOARCH=amd64 go build -o svg.linux
	GOOS=darwin GOARCH=amd64 go build -o svg.darwin
