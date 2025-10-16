# Database Documentation

## Overview

The application uses **PostgreSQL** as the primary data store. Database schema is managed using migrations located in `internal/infra/db/migrations/`.

## Connection Details

### Default Configuration
- **Host:** localhost (or `db` when using Docker)
- **Port:** 5432
- **User:** `POSTGRES_USER` (environment variable)
- **Password:** `POSTGRES_PASSWORD` (environment variable)
- **Database:** `POSTGRES_DB` (environment variable)
- **SSL Mode:** disabled

### Connection String
```
postgres://user:password@localhost:5432/database_name?sslmode=disable
```

---

## Schema

### Users Table

**Name:** `users`

**Purpose:** Stores application users with their basic information and role assignment.

#### Columns

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `uuid` | UUID | PRIMARY KEY | Unique user identifier |
| `email` | TEXT | NOT NULL | User's email address |
| `name` | TEXT | NOT NULL | User's full name |
| `role` | TEXT | NOT NULL | User role (attendee, host) |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT now() | Account creation timestamp |


## Migrations

Migrations are located in `internal/infra/db/migrations/` and use the naming convention:
```
XXXXXX_description.{up,down}.sql
```

### Current Migrations

#### Migration: 000001_create_users_table
- **Up:** Creates the `users` table
- **Down:** Drops the `users` table
- **File:** `000001_create_users_table.up.sql` / `000001_create_users_table.down.sql`
