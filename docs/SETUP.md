# Setup and Installation Guide

## Prerequisites

### Required
- **Go 1.24.0 or higher** (toolchain 1.24.6) - [Download Go](https://go.dev/dl/)
- **Docker and Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **PostgreSQL 15+** (or use Docker)

### Optional
- **Taskfile** - [Taskfile](https://taskfile.dev/) is a task runner for common operations. Install with:
  ```bash
  sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
  ```
  or 
  ```bash
    go install github.com/go-task/task/v3/cmd/task@latest
  ```

---

## Quick Start with Docker

### 1. Clone the Repository
```bash
git clone https://github.com/kapiw04/convenly.git
cd convenly
```

### 2. Create Environment File
Create a `.env` file in the project root:
```bash
touch .env
```
```env
POSTGRES_USER=convenly
POSTGRES_PASSWORD=convenly
POSTGRES_DB=convenly_db
```

### 3. Start Services
```bash
docker compose up -d
```

This starts:
- PostgreSQL database on port 5432
- Migrations to set up the database schema
- Go application on port 8080

### 4. Verify Setup
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status": "ok"}
```
---

## Common Commands

### Using Task (if installed)
```bash
task build         # Build the application
task run           # Run the application
task test          # Run tests
task migrate-up    # Run database migrations
task migrate-down  # Rollback migrations
```

### Using Go Directly
```bash
go build -o bin/app cmd/app/main.go
go run cmd/app/main.go
go test ./...
```

### Using Docker Compose
```bash
docker compose up -d       # Start services
docker compose logs -f     # View logs
docker compose down        # Stop services
docker compose ps          # List running services
docker compose exec db psql -U ${POSTGRES_USER} -d ${POSTGRES_DB}  # Access database
```

---

## Running Tests

### Unit Tests
Unit tests are located alongside the implementation files:
```bash
go test ./internal/infra/api -v
```

### Integration Tests
Integration tests use testcontainers to spin up real PostgreSQL instances:
```bash
go test ./test/integral -v
```

**Note:** Integration tests require Docker to be running.

### Generate Mocks
To regenerate mocks for testing:
```bash
go generate ./...
```

This will regenerate mocks using `mockgen` for interfaces marked with `//go:generate` directives.
