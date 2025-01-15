# skymates-api

## Getting Started

### Prerequisites

```bash
# List minimum requirements and dependencies
Go >= 1.22
PostgreSQL >= 15
```

### Installation

```bash
# Clone the repository
git clone git@github.com:jiyi27/skymates-api.git

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env file with your configuration

# Start the server
go run cmd/main.go
```

## API Documentation

### Authentication

The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints:

Obtain a JWT token by authenticating through the `/auth/login` endpoint
Include the token in the Authorization header of subsequent requests:

```
Authorization: Bearer <your_jwt_token>
```

Token Format:
```json
{
  "token": "eyJhbGciOiJ..."
}
```

### API Response Format:

All API responses follow a standard format:

```go
type Response struct {
  Message string      `json:"message"` // Response message
  Data    interface{} `json:"data"`    // Response payload
}
```

Example successful response:
```json
{
  "message": "Success",
  "data": {
    "id": 1,
    "name": "Example"
  }
}
```

Example error response:

```json
{
  "message": "Bad Request",
  "data": null
}
```


