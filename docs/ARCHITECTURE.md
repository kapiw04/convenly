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
│   ├── domain/                  # Domain layer (entities, interfaces)
│   │   └── user.go              # Core domain models
│   └── infra/                   # Infrastructure layer
│       ├── db/                  # Data access implementations
│       │   ├── postgresuser.go  # PostgreSQL repository
│       │   └── migrations/      # Database migrations
│       ├── http/                # HTTP server & routing
│       │   ├── handlers.go      # HTTP handlers
│       │   ├── response.go      # Response formatting
│       │   ├── router.go        # Route definitions
│       │   └── server.go        # Server lifecycle
│       └── log/                 # Logging
│           └── logger.go        # Centralized logging config
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
- Core business entities and interfaces
- Independent from infrastructure and framework code
- Defines contracts that infrastructure must implement (Repository Pattern)

### Infrastructure (`internal/infra/`)
- Technical implementations: database, HTTP routing, logging
- Implements domain interfaces (PostgresUserRepo implements UserRepo)
- Manages external integrations and framework specifics

## Technology Stack

- Go, Chi (HTTP router), PostgreSQL, Docker

See individual component files for detailed documentation.
