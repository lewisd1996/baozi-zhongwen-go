# Build stage
FROM golang:1.21 AS build-stage
WORKDIR /app
# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download
# Copy the rest of the source code
COPY . /app
# Build the application
# Update the build command to target the cmd directory
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /entrypoint ./cmd

# Deploy stage
FROM gcr.io/distroless/static-debian11 AS release-stage
WORKDIR /
# Copy the compiled binary and assets from the build stage
COPY --from=build-stage /entrypoint /entrypoint
COPY --from=build-stage /app/assets /assets
# Expose the port your app runs on
EXPOSE 3000
# Set the entry point to the compiled binary
ENTRYPOINT ["/entrypoint"]
