# Build stage
FROM golang:1.21 AS build-stage

# GET ARGUMENTS
ARG POSTGRES_USER
ARG POSTGRES_PASSWORD
ARG POSTGRES_HOST
ARG POSTGRES_PORT
ARG POSTGRES_DB

WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download
# Copy the rest of the source code
COPY . /app

# Run migrations and jet
RUN cd sql/migrations && goose postgres "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" up
RUN jet -dsn="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" -schema=public -path=./sql/.jet generate

# Build the application
# Update the build command to target the cmd directory
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint ./cmd

# Deploy stage
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
# Copy the compiled binary and assets from the build stage
COPY --from=build-stage /entrypoint /entrypoint
COPY --from=build-stage /app/assets /assets
# Expose the port your app runs on
EXPOSE 3000
# Use a non-root user for running the application
USER nonroot:nonroot
# Set the entry point to the compiled binary
ENTRYPOINT ["/entrypoint"]
