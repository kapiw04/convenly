# API Documentation

## Base URL
```
http://localhost:8080/api
```

## Endpoints

### Health Check

#### `GET /api/health`
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

#### `POST /api/register`
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
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Smith",
    "email": "alice@example.com",
    "password": "Secret123!"
  }'
```

---

### User Login

#### `POST /api/login`
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
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "alice@example.com",
  "name": "Alice Smith",
  "role": 0,
  "created_at": "2025-12-14T10:00:00Z"
}
```
**Status Code:** `200 OK`

**Response Headers:**
- `Set-Cookie: session-id=<session-token>; HttpOnly; Secure`

**Example cURL Request:**
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice@example.com",
    "password": "Secret123!"
  }'
```

---

### User Logout

#### `POST /api/logout`
Logs out the current user by invalidating their session token.

**Authentication Required:** Yes (via `session-id` cookie)

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
curl -X POST http://localhost:8080/api/logout \
  -H "Cookie: session-id=<session-token>"
```

---

### Get Current User Info

#### `GET /api/me`
Returns information about the currently authenticated user.

**Authentication Required:** Yes (via `session-id` cookie)

**Successful Response:**
```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "email": "alice@example.com",
  "name": "Alice Smith",
  "role": 0,
  "created_at": "2025-12-14T10:00:00Z"
}
```
**Status Code:** `200 OK`

**Role Values:**
| Value | Role |
|-------|------|
| 0 | Attendee |
| 1 | Host |

**Example cURL Request:**
```bash
curl -X GET http://localhost:8080/api/me \
  -H "Cookie: session-id=<session-token>"
```

---

### Become Host

#### `POST /api/become-host`
Promotes the current user from Attendee to Host role. Hosts can create events.

**Authentication Required:** Yes (via `session-id` cookie)

**Successful Response:**
```json
{
  "status": "ok"
}
```
**Status Code:** `200 OK`

**Example cURL Request:**
```bash
curl -X POST http://localhost:8080/api/become-host \
  -H "Cookie: session-id=<session-token>"
```

---

## Event Management

### Create Event

#### `POST /api/events/add`
Creates a new event. Only authenticated users with Host role can create events.

**Authentication Required:** Yes (via `session-id` cookie)
**Authorization Required:** Host role

**Request Body:**
```json
{
  "name": "Tech Conference 2025",
  "description": "Annual technology conference",
  "date": "2025-12-15T09:00:00Z",
  "latitude": 52.2297,
  "longitude": 21.0122,
  "fee": 99.99,
  "tags": ["technology", "networking"]
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
| `tags` | string[] | No | Array of tag names for the event |

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

Forbidden (user is not a Host):
```json
{
  "error": "forbidden"
}
```
**Status Code:** `403 Forbidden`

**Example cURL Request:**
```bash
curl -X POST http://localhost:8080/api/events/add \
  -H "Content-Type: application/json" \
  -H "Cookie: session-id=<session-token>" \
  -d '{
    "name": "Tech Conference 2025",
    "description": "Annual technology conference",
    "date": "2025-12-15T09:00:00Z",
    "latitude": 52.2297,
    "longitude": 21.0122,
    "fee": 99.99,
    "tags": ["technology"]
  }'
```

---

### List All Events

#### `GET /api/events`
Retrieves a list of all events with optional filtering and pagination.

**Authentication Required:** No

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `page` | int | No | Page number (starting from 1) |
| `page_size` | int | No | Number of items per page (1-100, default: 12) |
| `date_from` | string | No | Filter events from this date (RFC3339 or YYYY-MM-DD) |
| `date_to` | string | No | Filter events until this date (RFC3339 or YYYY-MM-DD) |
| `min_fee` | float | No | Minimum event fee |
| `max_fee` | float | No | Maximum event fee |
| `tags` | string | No | Comma-separated list of tag names |

**Successful Response:**
```json
[
  {
    "event_id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Tech Conference 2025",
    "description": "Annual technology conference",
    "date": "2025-12-15T09:00:00Z",
    "latitude": 52.2297,
    "longitude": 21.0122,
    "fee": 99.99,
    "organizer_id": "987fcdeb-51a2-43d7-9abc-123456789def",
    "tag": ["technology", "networking"]
  }
]
```
**Status Code:** `200 OK`

**Example cURL Requests:**
```bash
curl -X GET http://localhost:8080/api/events

curl -X GET "http://localhost:8080/api/events?page=1&page_size=10"

curl -X GET "http://localhost:8080/api/events?date_from=2025-01-01&max_fee=50&tags=music,outdoor"
```

---

### Get Event Details

#### `GET /api/events/{id}`
Retrieves detailed information about a specific event including attendee count.

**Authentication Required:** Yes (via `session-id` cookie)

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Event identifier |

**Successful Response:**
```json
{
  "event_id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Tech Conference 2025",
  "description": "Annual technology conference",
  "date": "2025-12-15T09:00:00Z",
  "latitude": 52.2297,
  "longitude": 21.0122,
  "fee": 99.99,
  "organizer_id": "987fcdeb-51a2-43d7-9abc-123456789def",
  "tag": ["technology"],
  "attendees_count": 42,
  "user_registered": true
}
```
**Status Code:** `200 OK`

**Response Fields:**
| Field | Type | Description |
|-------|------|-------------|
| `attendees_count` | int | Number of users registered for this event |
| `user_registered` | bool | Whether the current user is registered for this event |

**Example cURL Request:**
```bash
curl -X GET http://localhost:8080/api/events/123e4567-e89b-12d3-a456-426614174000 \
  -H "Cookie: session-id=<session-token>"
```

---

### Register for Event

#### `POST /api/events/{id}/register`
Registers the current user as an attendee for the specified event.

**Authentication Required:** Yes (via `session-id` cookie)

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Event identifier |

**Successful Response:**
```json
{
  "status": "ok"
}
```
**Status Code:** `200 OK`

**Example cURL Request:**
```bash
curl -X POST http://localhost:8080/api/events/123e4567-e89b-12d3-a456-426614174000/register \
  -H "Cookie: session-id=<session-token>"
```

---

### Unregister from Event

#### `DELETE /api/events/{id}/unregister`
Removes the current user's registration from the specified event.

**Authentication Required:** Yes (via `session-id` cookie)

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Event identifier |

**Successful Response:**
```json
{
  "status": "ok"
}
```
**Status Code:** `200 OK`

**Example cURL Request:**
```bash
curl -X DELETE http://localhost:8080/api/events/123e4567-e89b-12d3-a456-426614174000/unregister \
  -H "Cookie: session-id=<session-token>"
```

---

### Delete Event

#### `DELETE /api/events/{id}`
Deletes the specified event. Only the event organizer can delete their own events.

**Authentication Required:** Yes (via `session-id` cookie)
**Authorization Required:** Host role + Event owner

**URL Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | UUID | Event identifier |

**Successful Response:**
```json
{
  "status": "ok"
}
```
**Status Code:** `200 OK`

**Error Responses:**

Event not found:
```json
{
  "error": "event not found"
}
```
**Status Code:** `404 Not Found`

Not the event owner:
```json
{
  "error": "you can only delete your own events"
}
```
**Status Code:** `403 Forbidden`

**Example cURL Request:**
```bash
curl -X DELETE http://localhost:8080/api/events/123e4567-e89b-12d3-a456-426614174000 \
  -H "Cookie: session-id=<session-token>"
```

---

### Get My Events

#### `GET /api/my-events`
Retrieves events that the current user is hosting or attending.

**Authentication Required:** Yes (via `session-id` cookie)

**Successful Response:**
```json
{
  "hosting": [
    {
      "event_id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Tech Conference 2025",
      "description": "Annual technology conference",
      "date": "2025-12-15T09:00:00Z",
      "latitude": 52.2297,
      "longitude": 21.0122,
      "fee": 99.99,
      "organizer_id": "987fcdeb-51a2-43d7-9abc-123456789def"
    }
  ],
  "attending": [
    {
      "event_id": "456e7890-e89b-12d3-a456-426614174000",
      "name": "Music Festival",
      "description": "Summer music festival",
      "date": "2025-07-20T18:00:00Z",
      "latitude": 51.5074,
      "longitude": -0.1278,
      "fee": 150.00,
      "organizer_id": "111fcdeb-51a2-43d7-9abc-123456789def"
    }
  ]
}
```
**Status Code:** `200 OK`

**Example cURL Request:**
```bash
curl -X GET http://localhost:8080/api/my-events \
  -H "Cookie: session-id=<session-token>"
```
