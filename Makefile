.PHONY: build test clean ui-build lint

build: ui-build
	go build -o rampart ./cmd/rampart

test:
	go test -v ./...

clean:
	rm -f rampart
	rm -rf ui/dist

ui-build:
	@if [ -f ui/package.json ]; then cd ui && npm ci && npm run build; else echo "ui package.json not found, skipping ui-build"; fi

lint:
	go vet ./...
