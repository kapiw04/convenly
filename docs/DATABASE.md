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

## ERD Diagram
![erd-diagram](../assets/docs/convenly-db.png)
> Note: Project is not in Work in Progress stage. Not everything is yet implemented. 

## Schema

### Users Table

**Name:** `users`

**Purpose:** Stores application users with their basic information and role assignment.

#### Columns

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `user_id` | UUID | PRIMARY KEY | Unique user identifier |
| `email` | TEXT | NOT NULL | User's email address |
| `password_hash` | TEXT | NOT NULL | Hashed user password |
| `name` | TEXT | NOT NULL | User's full name |
| `role` | UUID | FOREIGN KEY REFERENCES roles(id) | User's role identifier |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT now() | Account creation timestamp |

### Role Table
**Name:** `roles`
**Purpose:** Defines user roles within the application.

#### Columns
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `role_id` | UUID | PRIMARY KEY | Unique role identifier |
| `name` | TEXT | UNIQUE, NOT NULL | Name of the role (e.g., Host, Attendee) |

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

#### Migration: 000002_create_roles_table
- **Up:** Creates the `roles` table
- **Down:** Drops the `roles` table
- **File:** `000002_create_roles_table.up.sql` / `000002_create_roles_table.down.sql`

#### Migration: 000003_add_role_to_user
- **Up:** Adds foreign key constraint to `users.role` referencing `roles.role_id` with `ON DELETE SET NULL`
- **Down:** Removes foreign key constraint from `users.role`
- **File:** `000003_add_role_to_user.up.sql` / `000003_add_role_to_user.down.sql`