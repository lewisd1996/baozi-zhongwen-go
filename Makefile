run: 
	@templ generate
	@cd internal/tailwind && npm run build
	@go run cmd/main.go

build:
	@templ generate
	@go build -o main cmd/main.go