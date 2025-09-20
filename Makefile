build:
	go build -o bin/aliasgen ./cmd/aliasgen
test:
	go test ./...
release:
	goreleaser release --clean