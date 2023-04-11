GOENV:=GOOS=linux GOARCH=amd64 CGO_ENABLED=0

.PHONY: build
build: build/main

.PHONY: package
package: main.zip

build/main: go.mod *.go
	mkdir -p build
	$(GOENV) go build -o build/ main.go

main.zip: build/main
	cd build && zip main.zip main
