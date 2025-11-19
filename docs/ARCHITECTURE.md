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
│   │   ├── userservice.go       # User use case orchestration
│   │   └── eventservice.go      # Event use case orchestration
│   ├── config/                  # Configuration management
│   ├── domain/                  # Domain layer (entities, interfaces, value objects)
│   │   ├── event/               # Event domain
│   │   │   └── event.go         # Event entity and repository interface
│   │   ├── security/            # Security-related domain contracts
│   │   │   ├── hasher.go        # Password hashing interface
│   │   │   └── mocks/           # Generated mocks for testing
│   │   │       └── mock_hasher.go
│   │   └── user/                # User domain
│   │       ├── user.go          # User entity and repository interface
│   │       ├── session.go       # Session entity and repository interface
│   │       ├── email.go         # Email value object with validation
│   │       ├── password.go      # Password value object with validation
│   │       ├── errors.go        # Domain-specific errors
│   │       └── mocks/           # Generated mocks for testing
│   │           ├── mock_userrepo.go
│   │           └── mock_sessionrepo.go
│   └── infra/                   # Infrastructure layer
│       ├── webapi/              # HTTP API
│       │   ├── handlers.go      # HTTP handlers (user & event endpoints)
│       │   ├── handlers_test.go # Unit tests for handlers
│       │   ├── middlewares.go   # Auth and ACL middleware
│       │   ├── request.go       # Request DTOs
│       │   ├── response.go      # Response formatting
│       │   ├── router.go        # Route definitions
│       │   └── server.go        # Server lifecycle
│       ├── db/                  # Data access implementations
│       │   ├── postgresuser.go     # PostgreSQL user repository
│       │   ├── postgressession.go  # PostgreSQL session repository
│       │   ├── postgresevent.go    # PostgreSQL event repository
│       │   └── migrations/         # Database migrations (9 migrations)
│       ├── log/                 # Logging
│       │   └── logger.go        # Centralized logging config
│       └── security/            # Security implementations
│           └── bcrypthasher.go  # Bcrypt password hasher
├── test/
│   └── integral/                # Integration tests
│       ├── pg.go                # PostgreSQL test container setup
│       ├── fixtures.go          # Test setup helpers
│       ├── http.go              # HTTP test helpers
│       ├── tx.go                # Transaction test helpers
│       ├── registeruser_test.go # User registration integration test
│       ├── login_test.go        # User login integration test
│       └── create_event_test.go # Event creation integration test
└── go.mod
```

## Layers

### Entry Point (`cmd/`)
- Application entry point and initialization
- Wires together all dependencies using dependency injection

### Application (`internal/app/`)
- Business logic and use cases orchestration
- **UserService**: Handles user registration, login, logout, and session management
- **EventService**: Handles event creation and retrieval
- Services depend on domain interfaces for data access

### Domain (`internal/domain/`)
- Core business entities, value objects, and interfaces
- Independent from infrastructure and framework code
- Defines contracts that infrastructure must implement (Repository Pattern)
- **User Domain**: User entity, Email and Password value objects, validation rules, Session management
- **Event Domain**: Event entity with location and organizer information
- **Security Domain**: Password hashing contracts

### Infrastructure (`internal/infra/`)
- Technical implementations: database, HTTP routing, logging, security
- **Database Layer**: PostgreSQL implementations for User, Session, and Event repositories
- **Web API Layer**: HTTP handlers, routing, middleware (authentication, ACL)
- **Security Layer**: Bcrypt password hashing implementation
- Implements domain interfaces (e.g., PostgresUserRepo implements UserRepo)
- Manages external integrations and framework specifics

## Authentication & Authorization

### Authentication Middleware
- Validates session cookies on protected routes
- Extracts user context from session and adds to request context
- Returns 401 Unauthorized if session is invalid or missing

### ACL Middleware
- Checks user roles against required permissions
- Built on top of authentication middleware
- Supports role-based access control (e.g., Host, Attendee roles)

### Protected Routes
- `GET /events` - Requires valid session
- `POST /events/add` - Requires valid session (user becomes organizer)

### Testing (`test/`)
- **Unit Tests**: Located alongside implementation files (e.g., `handlers_test.go`)
- **Integration Tests**: In `test/integral/` directory

## Technology Stack

- Go, Chi (HTTP router), PostgreSQL, Docker

See individual component files for detailed documentation.
