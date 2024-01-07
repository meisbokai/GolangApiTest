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