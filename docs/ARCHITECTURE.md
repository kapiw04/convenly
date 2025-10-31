# Convenly - Project Architecture

## Overview

Follows the standard Go project layout with clean layered architecture.

## Project Structure

```
.
├── cmd/                         # Entry point
│   └── app/
│       └── main.go              # Initializes app and wires dependencies
├── internal/
│   ├── app/                     # Business logic layer
│   │   └── userservice.go       # Use case orchestration
│   ├── config/                  # Configuration management
│   ├── domain/                  # Domain layer (entities, interfaces, value objects)
│   │   ├── security/            # Security-related domain contracts
│   │   │   └── hasher.go        # Password hashing interface
│   │   └── user/                # User domain
│   │       ├── user.go          # User entity and repository interface
│   │       ├── email.go         # Email value object with validation
│   │       ├── password.go      # Password value object with validation
│   │       ├── errors.go        # Domain-specific errors
│   │       └── mocks/           # Generated mocks for testing
│   │           └── mock_userrepo.go
│   └── infra/                   # Infrastructure layer
│       ├── api/                 # HTTP API
│       │   ├── handlers.go      # HTTP handlers
│       │   ├── handlers_test.go # Unit tests for handlers
│       │   ├── request.go       # Request DTOs
│       │   ├── response.go      # Response formatting
│       │   ├── router.go        # Route definitions
│       │   ├── server.go        # Server lifecycle
│       │   └── mocks/           # Generated mocks for testing
│       ├── db/                  # Data access implementations
│       │   ├── postgresuser.go  # PostgreSQL repository
│       │   └── migrations/      # Database migrations (8 migrations)
│       ├── log/                 # Logging
│       │   └── logger.go        # Centralized logging config
│       └── security/            # Security implementations
│           └── bcrypthasher.go  # Bcrypt password hasher
├── test/
│   └── integral/                # Integration tests
│       ├── pg.go                # PostgreSQL test container setup
│       ├── registeruser_test.go # User registration integration test
│       └── tx.go                # Transaction test helpers
└── go.mod
```

## Layers

### Entry Point (`cmd/`)
- Application entry point and initialization
- Wires together all dependencies using dependency injection

### Application (`internal/app/`)
- Business logic and use cases orchestration
- Services depend on domain interfaces for data access

### Domain (`internal/domain/`)
- Core business entities, value objects, and interfaces
- Independent from infrastructure and framework code
- Defines contracts that infrastructure must implement (Repository Pattern)

### Infrastructure (`internal/infra/`)
- Technical implementations: database, HTTP routing, logging
- Implements domain interfaces (PostgresUserRepo implements UserRepo)
- Manages external integrations and framework specifics

### Testing (`test/`)
- **Unit Tests**: Located alongside implementation files (e.g., `handlers_test.go`)
- **Integration Tests**: In `test/integral/` directory

## Technology Stack

- Go, Chi (HTTP router), PostgreSQL, Docker

See individual component files for detailed documentation.
