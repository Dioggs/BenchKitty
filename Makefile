BINARY := BenchKitty
PKG := ./...

.PHONY: run build test fmt vet check clean

run:
	go run . -r 10 -d 250 -o ./ -t GET -p https://example.com

build:
	go build -o $(BINARY) .

test:
	go test $(PKG)

fmt:
	gofmt -w .

vet:
	go vet $(PKG)

check: fmt vet test

clean:
	rm -f $(BINARY) benchmark.xlsx
