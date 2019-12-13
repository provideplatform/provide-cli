.PHONY: build clean install mod test

clean:
	rm -rf ./.bin 2>/dev/null || true
	rm -rf ./vendor 2>/dev/null || true
	rm ./prvd 2>/dev/null || true
	go fix ./...
	go clean -i

build: clean mod
	go fmt ./...
	go build -v -o ./.bin/prvd .

install: clean
	go install ./...

mod:
	go mod init 2>/dev/null || true
	go mod tidy
	go mod vendor

test: build
	# TODO
