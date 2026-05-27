.PHONY: fmt test bench web serve

GO_BUILD_CACHE ?= $(CURDIR)/.cache/go-build

fmt:
	gofmt -w cmd pkg

test:
	GOCACHE=$(GO_BUILD_CACHE) go test ./...

bench:
	GOCACHE=$(GO_BUILD_CACHE) go run ./cmd/opencsv bench -rows 100000 -cols 16 -out /tmp/opencsv-bench.csv

web:
	python3 -m http.server 43255 -d web

serve:
	GOCACHE=$(GO_BUILD_CACHE) go run ./cmd/opencsv serve -addr 127.0.0.1:43255 -web web
