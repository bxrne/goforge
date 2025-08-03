# goforge

A powerful CLI tool for building Go projects faster with pre-configured templates and integrations.
Scaffold complete projects with REST APIs, gRPC services, CLI applications, and best practices built-in.

Inspired by [kirimase for Next.JS](https://kirimase.dev) by @nicoalbanese.

## Features

### Core Commands
- `init`: Initialize a new Go project with a standard structure
- `add`: Add new packages and integrations to existing projects  
- `generate`: Scaffold resources like models, controllers, and services

### Built-in Templates & Integrations
- **REST API**: Complete API server with Gin, Swagger/OpenAPI documentation, and Prometheus metrics
- **gRPC Services**: Protocol buffer definitions, server/client implementations, and reflection
- **CLI Applications**: Cobra-based command-line tools with subcommands
- **Docker**: Production-ready containers for development and deployment
- **Database Support**: SQLite, PostgreSQL, MongoDB, MySQL, Oracle
- **Authentication**: JWT tokens, OAuth2 integration

### Development Tools
- **Metrics & Monitoring**: Prometheus metrics collection
- **API Documentation**: Auto-generated Swagger/OpenAPI specs
- **Code Generation**: Protocol buffers, database models, API handlers
- **Testing**: Test templates and coverage reporting
- **Linting & Formatting**: Pre-configured tools and standards

## Installation

```bash
go install github.com/bxrne/goforge@latest
```

## Quick Start

### Create a new API project
```bash
goforge init my-api --template api --db postgres --auth jwt
```

### Create a gRPC service
```bash
goforge init my-service --template grpc
```

### Create a CLI tool
```bash
goforge init my-cli --template cli
```

## Usage

```bash
goforge [command] [flags]
```

### Available Commands

#### Initialize a new project
```bash
goforge init <project-name> [flags]
```

**Flags:**
- `-t, --template`: Project template (api, cli, grpc) [default: api]
- `-d, --db`: Database type (sqlite, postgres, mongodb, mysql, oracle)
- `-a, --auth`: Authentication type (jwt, oauth2)
- `-g, --go-version`: Go version to use (auto-detected if not specified)

#### Add integrations to existing project
```bash
goforge add <package>
```

#### Generate scaffolded resources
```bash
goforge generate <resource>
```

## Project Structure

goforge creates projects with the following structure:

```
my-project/
├── cmd/
│   └── my-project/          # Application entrypoints
│       └── main.go
├── internal/                # Private application code
│   ├── api/                # REST API server
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── middleware/     # HTTP middleware
│   │   └── routes/         # Route definitions
│   ├── grpc/              # gRPC services
│   │   ├── proto/         # Protocol buffer definitions
│   │   ├── server/        # gRPC server implementation
│   │   └── client/        # gRPC client implementation
│   └── cli/               # CLI commands
│       └── commands/      # Command implementations
├── pkg/                   # Public library code
├── configs/               # Configuration files
├── Dockerfile            # Container definition
├── Makefile              # Build automation
└── go.mod                # Go module definition
```

## Development

### Prerequisites
- Go 1.21 or later
- Protocol Buffers compiler (protoc)
- Docker (optional)

### Setup Development Environment
```bash
make setup
```

### Build the project
```bash
make build
```

### Run tests
```bash
make test
```

### Generate protobuf code
```bash
make proto
```

### Generate API documentation
```bash
make swagger
```

### Available Make Targets
- `build`: Build the CLI binary
- `build-all`: Build for multiple platforms
- `install`: Install binary to GOPATH/bin
- `test`: Run tests
- `test-coverage`: Run tests with coverage
- `proto`: Generate protobuf code
- `swagger`: Generate Swagger documentation
- `lint`: Run linting
- `fmt`: Format code
- `tidy`: Tidy dependencies
- `clean`: Clean build artifacts
- `dev`: Run in development mode
- `docker-build`: Build Docker image
- `docker-run`: Run Docker container
- `setup`: Setup development environment

## Examples

### REST API with Database
```bash
goforge init blog-api --template api --db postgres --auth jwt
cd blog-api
make build
./bin/blog-api
```

### gRPC Microservice
```bash
goforge init user-service --template grpc
cd user-service
make proto
make build
```

### CLI Tool
```bash
goforge init my-tool --template cli
cd my-tool
make build
./bin/my-tool --help
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
