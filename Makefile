.PHONY: build clean install test

clean:
	rm -rf ./.bin 2>/dev/null || true
	rm ./provide-cli 2>/dev/null || true
	go fix ./...
	go clean -i

build: clean
	go fmt ./...
	go build -v -o ./.bin/provide-cli .

install: clean
	go install ./...

test: build
	# TODO
