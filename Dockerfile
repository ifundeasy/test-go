# Use the official Golang image as the base image
FROM golang:1.22.5-alpine AS builder

# Set environment variables for Go
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Create a working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o /test-go ./cmd/api/main.go

# Use a minimal image for the runtime
FROM alpine:latest

# Create a working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /test-go .

# Copy the .env file into the container
COPY .env .env

# Expose the port the app runs on
EXPOSE 8080

# Run the application
CMD ["./test-go"]
