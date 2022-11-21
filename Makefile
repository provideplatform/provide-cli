.PHONY: build build-binaries clean install mod test

clean:
	rm -rf ./.bin 2>/dev/null || true
	rm -rf ./vendor 2>/dev/null || true
	rm ./prvd 2>/dev/null || true
	go fix ./...
	go clean -i

build: clean mod
	go fmt ./...
	CGO_CFLAGS=-Wno-undef-prefix go build -v -o ./.bin/prvd ./cmd/prvd
	CGO_CFLAGS=-Wno-undef-prefix go build -v -o ./.bin/prvdnetwork ./cmd/prvdnetwork

build-binaries: clean mod
	go fmt ./...
	GOOS=darwin GOARCH=amd64 go build -o bin/prvd-amd64-darwin cmd/prvd/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/prvd-amd64-windows cmd/prvd/main.go
	GOOS=linux GOARCH=amd64 go build -o bin/prvd-amd64-linux cmd/prvd/main.go

install: build
	mkdir -p "${GOPATH}/bin"
	mv ./.bin/prvd "${GOPATH}/bin/prvd"
	mv ./.bin/prvdnetwork "${GOPATH}/bin/prvdnetwork"
	rm -rf ./.bin

mod:
	go mod init 2>/dev/null || true
	go mod tidy
	go mod vendor

test: build
	# TODO
