BINARY := gsd
VERSION := 0.1.0

.PHONY: build test lint clean install

build:
	go build -ldflags "-s -w" -o $(BINARY) .

test:
	go test ./... -v

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

install: build
	install -d ~/.local/bin
	install -m 755 $(BINARY) ~/.local/bin/$(BINARY)
