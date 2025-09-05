# Go Auth Service API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
All protected endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Endpoints

### Authentication

#### Register User
```http
POST /v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure_password123"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com"
  }
}
```

**Error Response (400 Bad Request):**
```json
{
  "success": false,
  "error": "Invalid email format"
}
```

#### Login User
```http
POST /v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "secure_password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "550e8400-e29b-41d4-a716-446655440000",
    "expires_in": 900,
    "token_type": "Bearer"
  }
}
```

#### Refresh Token
```http
POST /v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Logout
```http
POST /v1/auth/logout
Authorization: Bearer <access-token>
```

### User Management (Protected)

#### Get Current User
```http
GET /v1/auth/me
Authorization: Bearer <access-token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "email_verified": false,
    "created_at": "2024-01-15T10:30:00Z",
    "roles": ["user"]
  }
}
```

#### Update User Profile
```http
PUT /v1/auth/me
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "email": "newemail@example.com"
}
```

#### Change Password
```http
PUT /v1/auth/password
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "current_password": "old_password",
  "new_password": "new_secure_password123"
}
```

### Health & Monitoring

#### Health Check
```http
GET /health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0",
  "uptime": "2h30m15s"
}
```

#### Readiness Probe
```http
GET /ready
```

**Response (200 OK):**
```json
{
  "status": "ready",
  "service": "auth-service",
  "version": "1.0.0"
}
```

#### Liveness Probe
```http
GET /live
```

**Response (200 OK):**
```json
{
  "status": "alive"
}
```

#### Prometheus Metrics
```http
GET /metrics
```

Returns Prometheus-formatted metrics for monitoring.

## Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid input data |
| 401 | Unauthorized - Invalid or missing token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource doesn't exist |
| 409 | Conflict - Resource already exists |
| 422 | Unprocessable Entity - Validation failed |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Server error |

## Rate Limiting

- **Default Limit**: 60 requests per minute per IP
- **Burst Limit**: 10 requests
- **Headers**: Rate limit information is included in response headers:
  - `X-RateLimit-Limit`: Maximum requests per window
  - `X-RateLimit-Remaining`: Remaining requests in current window
  - `X-RateLimit-Reset`: Time when the rate limit resets

## CORS Policy

The service supports CORS for web applications:
- **Allowed Origins**: Configurable (default: localhost:3000, localhost:8080)
- **Allowed Methods**: GET, POST, PUT, DELETE, OPTIONS, PATCH
- **Allowed Headers**: Accept, Authorization, Content-Type, X-CSRF-Token, etc.
- **Credentials**: Supported when configured

## Security Headers

All responses include security headers:
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Content-Security-Policy: [configured policy]`
- `Referrer-Policy: strict-origin-when-cross-origin`

## Example Usage

### JavaScript (Fetch API)
```javascript
// Register user
const registerUser = async (email, password) => {
  const response = await fetch('http://localhost:8080/v1/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });
  
  return response.json();
};

// Login user
const loginUser = async (email, password) => {
  const response = await fetch('http://localhost:8080/v1/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });
  
  const data = await response.json();
  if (data.success) {
    localStorage.setItem('access_token', data.data.access_token);
    localStorage.setItem('refresh_token', data.data.refresh_token);
  }
  
  return data;
};

// Make authenticated request
const getProfile = async () => {
  const token = localStorage.getItem('access_token');
  const response = await fetch('http://localhost:8080/v1/auth/me', {
    headers: {
      'Authorization': `Bearer ${token}`,
    },
  });
  
  return response.json();
};
```

### cURL Examples
```bash
# Register user
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login user
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Get user profile (replace TOKEN with actual token)
curl -X GET http://localhost:8080/v1/auth/me \
  -H "Authorization: Bearer TOKEN"

# Health check
curl -X GET http://localhost:8080/health
```

## Postman Collection

Import the following collection for easy API testing:

```json
{
  "info": {
    "name": "Go Auth Service",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080"
    },
    {
      "key": "access_token",
      "value": ""
    }
  ],
  "item": [
    {
      "name": "Register User",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"email\": \"test@example.com\",\n  \"password\": \"password123\"\n}"
        },
        "url": {
          "raw": "{{base_url}}/v1/auth/register",
          "host": ["{{base_url}}"],
          "path": ["v1", "auth", "register"]
        }
      }
    },
    {
      "name": "Health Check",
      "request": {
        "method": "GET",
        "url": {
          "raw": "{{base_url}}/health",
          "host": ["{{base_url}}"],
          "path": ["health"]
        }
      }
    }
  ]
}
```
