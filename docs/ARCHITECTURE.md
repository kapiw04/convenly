# Convenly - Project Architecture

## Layers

### Entry Point (`cmd/`)
- Application entry point and initialization
- Wires together all dependencies using dependency injection

### Application (`internal/app/`)
- Business logic and use cases orchestration
- **UserService**: Handles user registration, login, logout, session management, and role promotion
- **EventService**: Handles event CRUD, filtering, attendance registration, and organizer-specific queries
- Services depend on domain interfaces for data access

### Domain (`internal/domain/`)
- Core business entities, value objects, and interfaces
- Independent from infrastructure and framework code
- Defines contracts that infrastructure must implement (Repository Pattern)
- **User Domain**: User entity, Email and Password value objects, validation rules, Session management, Roles (Attendee, Host)
- **Event Domain**: Event entity with location, organizer, tags, and filtering capabilities
- **Security Domain**: Password hashing contracts

### Infrastructure (`internal/infra/`)
- Technical implementations: database, HTTP routing, logging, security
- **Database Layer**: PostgreSQL implementations for User, Session, Event, and Tag repositories
- **Web API Layer**: HTTP handlers, routing, middleware (authentication, ACL), CORS configuration
- **Security Layer**: Bcrypt password hashing implementation
- Implements domain interfaces (e.g., PostgresUserRepo implements UserRepo)
- Manages external integrations and framework specifics

### Frontend (`frontend/`)
- SvelteKit-based single-page application
- Component library using shadcn-svelte UI components
- User authentication state management via Svelte stores
- Event browsing, creation, and registration interfaces

## Authentication & Authorization

### Authentication Middleware
- Validates session cookies (`session-id`) on protected routes
- Extracts user context from session and adds to request context
- Returns 401 Unauthorized if session is invalid or missing

### ACL Middleware
- Checks user roles against required permissions
- Built on top of authentication middleware
- Supports role-based access control:
  - **Attendee (role=0)**: Can browse events, register for events, view their registrations
  - **Host (role=1)**: All Attendee permissions + can create and delete their own events

## Technology Stack

### Backend
- **Go 1.24** - Primary backend language
- **Chi** - HTTP router with middleware support
- **PostgreSQL** - Primary database
- **Docker & Docker Compose** - Containerization and orchestration
- **Taskfile** - Task runner for development commands

### Frontend
- **SvelteKit** - Full-stack web framework
- **TypeScript** - Type-safe JavaScript
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Utility-first CSS framework
- **shadcn-svelte** - UI component library

### Testing
- **gotestsum** - Go test runner
- **testcontainers** - Integration testing with real PostgreSQL
- **mockgen** - Mock generation for unit tests
