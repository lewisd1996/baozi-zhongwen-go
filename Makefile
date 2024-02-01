run: 
	@templ generate
	@cd internal/tailwind && npm run build
	@go run cmd/main.go

build:
	@go build -o main cmd/main.go