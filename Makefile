start:
	@export APP_ENV=development && air

run: 
	@templ generate
	@cd internal/tailwind && npm run build
	@go run cmd/main.go

create-migration:
	@cd internal/sql/migrations && goose create $(name) sql

migrate-up:
	@source .env && cd internal/sql/migrations && goose postgres "postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable" up

migrate-down:
	@source .env && cd internal/sql/migrations && goose postgres "postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable" down

migrate-status:
	@source .env && cd internal/sql/migrations && goose postgres "postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable" status

jet:
	@source .env && jet -dsn="postgres://$$POSTGRES_USER:$$POSTGRES_PASSWORD@$$POSTGRES_HOST:$$POSTGRES_PORT/$$POSTGRES_DB?sslmode=disable" -schema=public -path=./internal/sql/.jet generate