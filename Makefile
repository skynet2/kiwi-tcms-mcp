.PHONY: build test test-e2e lint generate run tidy clean

BINARY := bin/kiwi-tcms-mcp

build:
	go build -o $(BINARY) ./cmd/kiwi-tcms-mcp

test:
	go test -p 1 -timeout 60s ./...

test-e2e:
	go test -tags=e2e -p 1 -timeout 120s -count=1 -v ./...

lint:
	golangci-lint run

generate:
	go generate ./...

run:
	go run ./cmd/kiwi-tcms-mcp serve

tidy:
	go mod tidy

clean:
	rm -rf bin/
