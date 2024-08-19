# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=test-go
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER=docker
DOCKER_COMPOSE=docker-compose
DOCKER_IMAGE=test-go:latest

# Proto file
PROTO_DIR=proto
PROTO_OUT_DIR=internal/adapters/primary/grpc
PROTOC=$(GOCMD) run github.com/golang/protobuf/protoc-gen-go
PROTOC_GEN_GO=protoc-gen-go
PROTOC_GEN_GO_GRPC=protoc-gen-go-grpc

# Build the project
all: test build

build: 
  swag init -g cmd/http/server.go -o internal/adapters/primary/http/swagger
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

# Generate protobuf files
proto:
	protoc -I=$(PROTO_DIR) --go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative --go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative $(PROTO_DIR)/*.proto
	
# Run the project
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go
	./$(BINARY_NAME)

# Test the project
test:
	$(GOTEST) -v ./...

# Lint the project using golangci-lint
lint:
	$(GOCMD) run github.com/golangci/golangci-lint/cmd/golangci-lint run

# Clean the build files
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run Docker Compose
docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

# Build Docker image
docker-build:
	$(DOCKER) build -t $(DOCKER_IMAGE) .

# Tidy Go modules
tidy:
	$(GOMOD) tidy

# Update Go modules
mod-update:
	$(GOMOD) get -u ./...

# Help command to list available make targets
help:
	@echo "Makefile commands:"
	@echo "  make build       - Build the Go project"
	@echo "  make proto       - Generate protobuf files"
	@echo "  make run         - Build and run the project"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Run golangci-lint"
	@echo "  make clean       - Clean build files"
	@echo "  make docker-up   - Start services using Docker Compose"
	@echo "  make docker-down - Stop services using Docker Compose"
	@echo "  make docker-build- Build Docker image"
	@echo "  make tidy        - Tidy up Go modules"
	@echo "  make mod-update  - Update Go modules to the latest version"
	@echo "  make help        - Display this help message"
