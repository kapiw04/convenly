# Convenly - Project Architecture

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

## Technology Stack

- Go, Chi (HTTP router), PostgreSQL, Docker

See individual component files for detailed documentation.
