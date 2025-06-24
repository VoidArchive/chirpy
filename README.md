# Chirpy

Chirpy is a simple social media API that allows users to post short messages (chirps), manage user accounts, and interact with a Twitter-like platform.

## Features

- **User Management**: User registration, authentication, and profile updates
- **Chirp Management**: Create, read, and delete short messages
- **Authentication**: JWT-based authentication with refresh tokens
- **Premium Membership**: Chirpy Red subscription via webhook integration
- **Filtering & Sorting**: Get chirps by author and sort by creation date

## API Endpoints

### Health Check
- `GET /api/healthz` - Health check endpoint

### User Management
- `POST /api/users` - Create a new user account
- `PUT /api/users` - Update user email and password (requires authentication)
- `POST /api/login` - Login and receive access/refresh tokens
- `POST /api/refresh` - Get new access token using refresh token
- `POST /api/revoke` - Revoke a refresh token

### Chirps
- `GET /api/chirps` - Get all chirps (supports filtering and sorting)
  - Query parameters:
    - `author_id` (optional) - Filter chirps by author UUID
    - `sort` (optional) - Sort order: `asc` (default) or `desc`
- `GET /api/chirps/{chirpID}` - Get a specific chirp by ID
- `POST /api/chirps` - Create a new chirp (requires authentication)
- `DELETE /api/chirps/{chirpID}` - Delete a chirp (requires authentication, author only)

### Webhooks
- `POST /api/polka/webhooks` - Webhook for Polka payment processing (requires API key)

### Admin
- `GET /admin/metrics` - View server metrics
- `POST /admin/reset` - Reset all users (development only)

## Authentication

Chirpy uses JWT-based authentication with access and refresh tokens:

1. **Access Tokens**: Short-lived (1 hour) tokens for API access
2. **Refresh Tokens**: Long-lived (60 days) tokens for obtaining new access tokens

### Using Authentication

1. Create an account with `POST /api/users`
2. Login with `POST /api/login` to receive tokens
3. Include access token in requests: `Authorization: Bearer <access_token>`
4. Use refresh token to get new access tokens when expired

## Request/Response Examples

### Create User
```bash
POST /api/users
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

### Login
```bash
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

Response:
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "56aa826d22baab4b5ec2cea41a59ecbba..."
}
```

### Create Chirp
```bash
POST /api/chirps
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "body": "Hello, world!"
}
```

### Get Chirps with Filtering
```bash
# Get all chirps (ascending order)
GET /api/chirps

# Get chirps by specific author (descending order)
GET /api/chirps?author_id=123e4567-e89b-12d3-a456-426614174000&sort=desc
```

## Setup and Installation

### Prerequisites
- Go 1.21+
- PostgreSQL
- Environment variables configured

### Environment Variables
Create a `.env` file with:
```env
DB_URL="postgres://username:password@localhost:5432/chirpy?sslmode=disable"
JWT_SECRET="your-jwt-secret-key"
POLKA_KEY="your-polka-api-key"
PLATFORM="dev"
```

### Database Setup
1. Create PostgreSQL database
2. Run migrations: `goose up`
3. Generate SQLC code: `sqlc generate`

### Running the Server
```bash
go build -o chirpy
./chirpy
```

The server will start on port 8080.

## Features in Detail

### Chirp Validation
- Maximum 140 characters
- Profanity filter replaces bad words with "****"
- Supported bad words: kerfuffle, sharbert, fornax

### Chirpy Red
- Premium membership feature
- Activated via Polka webhook integration
- Provides enhanced user status

### Security
- Passwords hashed with bcrypt
- JWT tokens with expiration
- API key authentication for webhooks
- Users can only modify their own content

## Development

### Project Structure
```
chirpy/
├── internal/
│   ├── auth/          # Authentication utilities
│   └── database/      # Database models and queries
├── sql/
│   ├── queries/       # SQL queries
│   └── schema/        # Database migrations
├── handler_*.go       # HTTP handlers
├── main.go           # Server entry point
└── README.md
```

### Testing
```bash
# Run auth package tests
go test ./internal/auth

# Build project
go build -o chirpy
```

## License

This project is part of a learning exercise and is not intended for production use.