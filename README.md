# Test-Go

Test-Go is a sample project demonstrating a hexagonal architecture implementation using the Go programming language. The project leverages Fiber as the web framework, gRPC for microservices communication, MongoDB for data persistence, Redis for caching, and RabbitMQ for message queuing.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Environment Configuration](#environment-configuration)
- [Generating Documentation](#generating-documentation)
- [Generating Proto Files](#generating-proto-files)
- [Running the Application](#running-the-application)
  - [Using Docker](#using-docker)
  - [Without Docker](#without-docker)
- [Running Tests](#running-tests)

## Features

- RESTful API built with Fiber
- gRPC support for microservices
- MongoDB for data persistence
- Redis for caching
- RabbitMQ for message queuing
- Swagger documentation generation
- Unit tests with mocking for HTTP, gRPC, MongoDB, and Redis

## Installation

To set up the project locally, follow these steps:

1. **Clone the repository:**

    ```bash
    git clone https://github.com/ifundeasy/test-go.git
    cd test-go
    ```

2. **Install Go modules:**

    Ensure you have Go installed on your machine. Run the following command to download and install the necessary dependencies:

    ```bash
    go mod tidy
    ```

## Environment Configuration

This project uses environment variables to manage configurations. These variables are loaded from a `.env` file.

1. **Create a `.env` file in the root of your project:**

    ```bash
    touch .env
    ```

2. **Add the necessary environment variables:**

    ```plaintext
    MONGO_URI=mongodb://localhost:27017/your_db
    REDIS_URI=redis://username:password@localhost:6379
    RABBITMQ_URI=amqp://guest:guest@localhost:5672/
    PORT=8080
    ```

## Generating Documentation

This project uses `swaggo` to generate API documentation in the OpenAPI format.

1. **Install swaggo:**

    If you haven't installed `swaggo` yet, do so by running:

    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```

2. **Generate Swagger documentation:**

    From the root of the project, run:

    ```bash
    swag init -g cmd/api/main.go --output internal/swagger
    ```

    The documentation will be generated in the `internal/swagger` directory.

## Generating Proto Files

For gRPC, you'll need to generate Go files from your `.proto` definitions:

1. **Install Protocol Buffers:**

    Follow the instructions from the [official Protocol Buffers documentation](https://grpc.io/docs/protoc-installation/) to install `protoc`.

2. **Install Go gRPC plugin:**

    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```

3. **Generate Go files from .proto:**

    Assuming your `.proto` files are in the `proto/` directory, run:

    ```bash
    protoc --go_out=. --go-grpc_out=. proto/*.proto
    ```

    This will generate the necessary Go files for gRPC in your project.

## Running the Application

### Using Docker

1. **Build the Docker image:**

    ```bash
    docker build -t test-go .
    ```

2. **Run the Docker container:**

    ```bash
    docker run --env-file .env -p 8080:8080 test-go
    ```

    The application will be accessible at `http://localhost:8080`.

### Without Docker

1. **Run the application:**

    ```bash
    go run cmd/api/main.go
    ```

    The API will be accessible at `http://localhost:8080`.

2. **Run the gRPC server:**

    ```bash
    go run cmd/grpc/server.go
    ```

    The gRPC server will listen on port `50051`.

## Running Tests

To ensure the application works as expected, you can run unit tests with mocking.

1. **Run the tests:**

    ```bash
    go test ./test/unit/...
    ```

    This command will run all the unit tests in the `test/unit` directory.

2. **Check test coverage:**

    If you want to check the test coverage, run:

    ```bash
    go test ./test/unit/... -coverprofile=coverage.out
    go tool cover -html=coverage.out
    ```

    This will generate a coverage report that you can view in your browser.

## Conclusion

This project provides a comprehensive example of using modern Go techniques and tools to build a scalable and maintainable application. With dotenv for environment management, Docker for containerization, and robust testing practices, this project serves as a solid foundation for building production-grade services. Feel free to explore and modify the code to suit your needs. Contributions are welcome!

