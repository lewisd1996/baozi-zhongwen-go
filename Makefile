run: 
	@templ generate
	@cd internal/tailwind && npm run build
	@go run cmd/main.go