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
