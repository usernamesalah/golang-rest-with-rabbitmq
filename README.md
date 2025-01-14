# Tenant Apps API

## Prerequesites

1. Docker
2. Golang-Migrate
3. Swag
4. Golang
5. RabbitMQ
6. Mockery

## Development Guide

### Getting Started

#### Running with docker
1. Run `make dev args="up"` to run local environment with hot reload, everything inside `args` is basically the `docker compose` arguments. You can exit by `ctrl+c`.
2. The db migration will automatically running.
3. If you want to tear down all the containers you can run `make dev args="down"` or add `-v` to the args if you want to remove the volume too, for example `make dev args="down -v"`
4. To rebuild the docker image if you change something, you can run `make dev args="build"`. This will left dangling images, you need to remove the dangling images manually.
5. dont forget to change value `config.yaml` file in root folder and change the environment variables to your own.

#### Running without docker
1. Run only webservice `go run cmd/main.go`
2. Tenant Service with cli
    - Create tenant `go run cmd/main.go [tenant-name]`
    - Process payload tenant `go run cmd/main.go [client-id] [tenant-payload]`
    - Delete tenant `go run cmd/main.go [client-id]`

#### RabbitMQ

When you run `make dev args="up"`, RabbitMQ will be automatically started as part of the local environment setup. This ensures that RabbitMQ is running and ready for use without any additional manual steps.

To verify that RabbitMQ is running, you can access the RabbitMQ management UI by visiting [http://localhost:15672](http://localhost:15672) with default credentials of `guest` and `guest`.

You can access the RabbitMQ management UI by visiting [http://localhost:15672](http://localhost:15672) with default credentials of `guest` and `guest`


### Create Migration

You need to install `golang-migrate` manually in your device.

You can install it by `go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest`.

Run `migrate create -ext sql -dir migration {migration_name}` to create your migration file

### Running DB Migration manually

Please make sure the `tenant-apps` and `database` containers are running.

everything inside `args` are the `golang-migrate` arguments

Run `make dev-migrate args="up"` to create or modify tables

Run `make dev-migrate args="drop -f"` to drop all tables

Run `make dev-migrate args="down [migration_version]"` to delete some tables based on migration version

Run `make dev-migrate args="force [migration_version]"` to fix the dirty migration

## Testing Guide

### Running Test

Run `make test`

### Create Mocks

Run `make mock` if using windows `make mock-win`

## Deployment Guide

### DB Migration

DB Migration on Staging or Prod need to do it manually by running `migrate -database "mysql://<DB_USER>:<DB_PASS>@tcp(<DB_HOST>:<DB_PORT>)/<DB_NAME>" -path migration [up|<any migrate argument>]`

### Deployment

The docker setup is inside `deployments/deploy` folder


## API Documentation
You need to install [swag](https://github.com/swaggo/swag) manually in your device.

You can install it by `go install github.com/swaggo/swag/cmd/swag`.

We use [swag](https://github.com/swaggo/swag) to generate necearry Swagger files for API documentation. Everytime we run `make gen-swagger`, the Swagger documentation will be updated.

To configure your API documentation, please refer to [swag's declarative comments format](https://github.com/swaggo/swag#declarative-comments-format) and [examples](https://github.com/swaggo/swag#examples).

To access the documentation, after you run the app please visit [API DOCUMENTATION](http://localhost:8080/v1/docs/index.html).


## Project Structure

- `cmd/webservice/main.go`: The entry point of the application where the server is initialized and started.
- `internal/api/http/router.go`: Configures all API routes and their corresponding handlers.
- `internal/container/container.go`: Manages dependency injection and application-wide component initialization.
- `internal/api/http/handler`: Contains HTTP handlers that process incoming requests and return responses.
- `internal/api/http/middleware`: Houses middleware functions for authentication, logging, error handling etc.
- `internal/api/http/request`: Defines request structs and validation for incoming API requests.
- `internal/api/http/response`: Contains response structs and helper functions for API responses.
- `internal/api/http/routes`: Groups related API routes and their configurations.
- `internal/model`: Contains domain models and entities that represent the core business objects and data structures.
- `internal/repository`: Data access layer that handles database operations and data persistence.
- `internal/service`: Implements core business logic and coordinates between different layers.
- `internal/usecase`: Orchestrates the flow of data and business rules between different services.
- `infrastructure/config`: Manages application configuration from environment variables and config files.
- `infrastructure/database`: Handles database connections, migrations and database-specific configurations.
- `infrastructure/database/migrations`: Contains database migration files for schema changes and data updates.
- `docs`: Stores auto-generated API documentation and other project documentation.
- `deployments`: Contains Docker files, deployment scripts and environment configurations.
- `pkg`: Reusable packages and utilities that can be shared across projects.
- `internal/test`: Contains test files, test utilities, and test configurations.

## Libraries to Use

- **pgx**: For database interaction with PostgreSQL. It provides a pure Go driver and tools for PostgreSQL, including support for prepared statements, transactions, and connection pooling.s
- **viper**: For configuration management. It supports reading from various configuration file formats and integrates seamlessly with environment variables and flags to provide a flexible configuration solution.
- **cobra**: For CLI integration. It is a library for creating powerful and modern command-line applications, supporting subcommands, flags, and automatic help generation.
- **logrus**: For structured logging. It provides a simple API for logging with leveled logs, hooks, and formatting options, making it easier to manage log output.
- **lumberjack**: For log file rotation. It handles rolling of log files based on size, date, and other criteria, helping to manage log file sizes and retention.
- **echo**: For building the REST API. It is a high performance, minimalist Go web framework that focuses on simplicity and developer productivity, with built-in middleware support.
- **amqp091**: For RabbitMQ communication. This library implements the AMQP 0.9.1 protocol, which is used for messaging between services and components in a microservices architecture.
- **docker**: For containerization. It provides a convenient way to package and run our application with its dependencies in a containerized environment.
- **golang-migrate**: For database migration management. It provides a simple way to manage database schema changes and data migrations
- **air**: For hot reloading. It is a live reload tool for Go applications, allowing developers to quickly see the effects of code changes without having to manually restart the application.
- **mockery**: For generating mock objects. It is a mock code autogenerator that supports a wide range of mocking libraries and test frameworks.
