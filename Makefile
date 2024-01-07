serve:
	go run cmd/api/main.go

generateDocumentation:
	swag init -g ./cmd/api/main.go --pd --parseDepth 3