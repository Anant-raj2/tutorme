default:
	@~/go/bin/templ generate
	@go run cmd/api/main.go
