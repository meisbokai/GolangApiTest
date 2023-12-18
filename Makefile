build:
	@go build -o bin/golangapitest

run: build
	@./bin/golangapitest

test:
	@go test -v ./...