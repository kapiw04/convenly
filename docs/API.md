# API Documentation

## Base URL
```
http://localhost:8080
```

## Endpoints

### Health Check

#### `GET /health`
Returns the health status of the API.

**Response:**
```json
{
  "status": "ok"
}
```

**Status Code:** `200 OK`

---

### User Registration

#### `POST /register`
Creates a new user account with the provided credentials.

**Request Body:**
```json
{
  "name": "Alice Smith",
  "email": "alice@example.com",
  "password": "Secret123!"
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | User's full name |
| `email` | string | Yes | User's email address (must be valid format) |
| `password` | string | Yes | User's password (see validation rules below) |

**Password Validation Rules:**
- Minimum length: 8 characters
- Maximum length: 20 characters
- Must contain at least one uppercase letter (A-Z)
- Must contain at least one lowercase letter (a-z)
- Must contain at least one digit (0-9)
- Must contain at least one special character (!@#~$%^&*()+ |_)

**Email Validation Rules:**
- Must be a valid email format
- Will be converted to lowercase
- Leading/trailing whitespace will be trimmed

**Example cURL Request:**
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith",
    "email": "alice@example.com",
    "password": "Secret123!"
  }'
```

---

### User Login

#### `POST /login`
Authenticates an existing user and returns a session token that can be used for future requests.

**Request Body:**
```json
{
  "email": "alice@example.com",
  "password": "Secret123!"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | Registered email (case-insensitive) |
| `password` | string | Yes | Password (same validation as registration) |


**Successful Response:**
```json
{
  "status": "ok"
}
```
**Status Code:** `200 OK`

**Response Headers:**
- `Set-Cookie: session_id=<session-token>; HttpOnly; Secure`

**Example cURL Request:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "Secret123!"
  }'
```

---

### User Logout

#### `POST /logout`
Logs out the current user by invalidating their session token.

**Authentication Required:** Yes (via `Authorization` header)

**Request Headers:**
| Header | Value | Required | Description |
|--------|-------|----------|-------------|
| `Authorization` | `<session-token>` | Yes | Session token obtained from login |

**Successful Response:**
```json
{
  "status": "ok"
}
```
**Status Code:** `200 OK`

**Error Response:**
```json
{
  "error": "missing session ID"
}
```
**Status Code:** `400 Bad Request`

**Example cURL Request:**
```bash
curl -X POST http://localhost:8080/logout \
  -H "Authorization: <session-token>"
```

---

## Event Management

### Create Event

#### `POST /events/add`
Creates a new event. Only authenticated users can create events.

**Authentication Required:** Yes (via `session_id` cookie)

**Request Body:**
```json
{
  "name": "Tech Conference 2025",
  "description": "Annual technology conference",
  "date": "2025-12-15T09:00:00Z",
  "latitude": 52.2297,
  "longitude": 21.0122,
  "fee": 99.99
}
```

**Request Fields:**
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Event name |
| `description` | string | Yes | Event description |
| `date` | string | Yes | Event date in ISO 8601 format (RFC3339) |
| `latitude` | float64 | Yes | Latitude coordinate of event location |
| `longitude` | float64 | Yes | Longitude coordinate of event location |
| `fee` | float32 | Yes | Event entrance fee |

**Successful Response:**
```json
{
  "status": "ok"
}
```
**Status Code:** `201 Created`

**Error Responses:**

Unauthorized (missing or invalid session):
```json
{
  "error": "unauthorized"
}
```
**Status Code:** `401 Unauthorized`

Bad request (invalid data):
```json
{
  "error": "bad request: <error details>"
}
```
**Status Code:** `400 Bad Request`

**Example cURL Request:**
```bash
curl -X POST http://localhost:8080/events/add \
  -H "Content-Type: application/json" \
  -H "Cookie: session_id=<session-token>" \
  -d '{
    "name": "Tech Conference 2025",
    "description": "Annual technology conference",
    "date": "2025-12-15T09:00:00Z",
    "latitude": 52.2297,
    "longitude": 21.0122,
    "fee": 99.99
  }'
```

---

### List All Events

#### `GET /events`
Retrieves a list of all events.

**Authentication Required:** Yes (via `session_id` cookie)

**Successful Response:**
```json
[
  {
    "EventID": "123e4567-e89b-12d3-a456-426614174000",
    "Name": "Tech Conference 2025",
    "Description": "Annual technology conference",
    "Date": "2025-12-15T09:00:00Z",
    "Latitude": 52.2297,
    "Longitude": 21.0122,
    "Fee": 99.99,
    "OrganizerID": "987fcdeb-51a2-43d7-9abc-123456789def"
  }
]
```
**Status Code:** `200 OK`

**Error Response:**

Unauthorized (missing or invalid session):
```json
{
  "error": "unauthorized"
}
```
**Status Code:** `401 Unauthorized`

**Example cURL Request:**
```bash
curl -X GET http://localhost:8080/events \
  -H "Cookie: session_id=<session-token>"
```
