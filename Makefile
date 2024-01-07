serve:
	go run cmd/api/main.go

generateDocumentation:
	swag init -g ./cmd/api/main.go --pd --parseDepth 3

seed:
	go run cmd/seed/main.go

mig-up:
	go run cmd/migration/main.go -up
mig-down:
	go run cmd/migration/main.go -down

# TODO: Fix. Command does not work when using `make`
# mock:
# 	docker run -v "$PWD":/src -w /src vektra/mockery --all --output=internal/mocks

coverage:
	go test -v ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html