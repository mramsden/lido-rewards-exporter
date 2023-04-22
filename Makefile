GOENV:=GOOS=linux GOARCH=amd64 CGO_ENABLED=0

.PHONY: build
build: build/main

build/main: go.mod *.go
	mkdir -p build
	$(GOENV) go build -o build/ main.go
