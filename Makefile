.PHONY: fmt test bench web

fmt:
	gofmt -w cmd pkg

test:
	go test ./...

bench:
	go run ./cmd/opencsv bench -rows 100000 -cols 16 -out /tmp/opencsv-bench.csv

web:
	python3 -m http.server 43255 -d web
