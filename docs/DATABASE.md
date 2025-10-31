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

## Schema

### Users Table

**Name:** `users`

**Purpose:** Stores application users with their basic information and role assignment.

#### Columns

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `user_id` | UUID | PRIMARY KEY | Unique user identifier |
| `email` | TEXT | NOT NULL | User's email address (stored in lowercase) |
| `password_hash` | TEXT | NOT NULL | Hashed user password using bcrypt |
| `name` | TEXT | NOT NULL | User's full name |
| `role` | SMALLINT | FOREIGN KEY REFERENCES roles(role_id) | User's role identifier |
| `created_at` | TIMESTAMPTZ | NOT NULL, DEFAULT now() | Account creation timestamp |

**Notes:**
- Passwords are hashed using bcrypt with minimum cost
- Email addresses are validated and normalized to lowercase
- Role is stored as an integer (0 = Attendee, 1 = Host)

### Role Table
**Name:** `roles`

**Purpose:** Defines user roles within the application.

#### Columns
| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `role_id` | SMALLINT | PRIMARY KEY | Unique role identifier (integer) |
| `name` | TEXT | UNIQUE, NOT NULL | Name of the role (e.g., Host, Attendee) |

**Pre-populated Roles:**
| role_id | name |
|---------|------|
| 0 | Attendee |
| 1 | Host |

**Notes:**
- Changed from UUID to SMALLINT for better performance and simplicity
- Role IDs are pre-defined integers for easy reference in application code

---

### Events Table

**Name:** `events`

**Purpose:** Stores event information created by organizers.

#### Columns

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `event_id` | UUID | PRIMARY KEY | Unique event identifier |
| `name` | TEXT | UNIQUE | Event name |
| `description` | TEXT | | Event description |
| `date` | DATE | | Event date |
| `geolocation` | POINT | | Geographic coordinates of the event location |
| `fee` | DECIMAL | | Event entrance fee |
| `organiser_id` | UUID | FOREIGN KEY REFERENCES users(user_id) ON DELETE CASCADE | User ID of the event organizer |

**Notes:**
- Events are deleted when the organizer (user) is deleted (CASCADE)
- Geolocation uses PostgreSQL POINT type for storing coordinates
- Event names must be unique across the system

---

### Attendances Table

**Name:** `atttendances` (note: typo in migration)

**Purpose:** Tracks which users are attending which events (many-to-many relationship).

#### Columns

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `event_id` | UUID | FOREIGN KEY REFERENCES events(event_id), PRIMARY KEY | Event identifier |
| `user_id` | UUID | FOREIGN KEY REFERENCES users(user_id), PRIMARY KEY | User identifier |

**Notes:**
- Composite primary key ensures a user can only register once per event
- Junction table implementing many-to-many relationship between users and events

---

### Tags Table

**Name:** `tags`

**Purpose:** Stores tag definitions for event categorization.

#### Columns

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `tag_id` | UUID | PRIMARY KEY | Unique tag identifier |
| `name` | TEXT | UNIQUE, NOT NULL | Tag name |

**Notes:**
- Tag names must be unique
- Used for categorizing and filtering events

---

### Event Tags Table

**Name:** `event_tag`

**Purpose:** Associates tags with events (many-to-many relationship).

#### Columns

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `event_id` | UUID | FOREIGN KEY REFERENCES events(event_id), PRIMARY KEY | Event identifier |
| `tag_id` | UUID | FOREIGN KEY REFERENCES tags(tag_id), PRIMARY KEY | Tag identifier |

**Notes:**
- Composite primary key ensures a tag can only be applied once per event
- Junction table implementing many-to-many relationship between events and tags
- Allows events to have multiple tags and tags to be used on multiple events

---

## Migrations

Migrations are located in `internal/infra/db/migrations/` and use the naming convention:
```
XXXXXX_description.{up,down}.sql
```