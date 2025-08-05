# Swagger API Documentation

This project includes comprehensive Swagger API documentation for all endpoints.

## Accessing Swagger UI

Once the server is running, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Available Endpoints

### Authentication Endpoints (`/api/v1/auth`)
- `POST /v1/auth/register` - Register a new user
- `POST /v1/auth/login` - Login user
- `GET /v1/auth/profile` - Get user profile (requires authentication)

### User Management Endpoints (`/api/v1/users`)
- `GET /v1/users` - Get all users with pagination and filtering
- `POST /v1/users` - Create a new user
- `GET /v1/users/{id}` - Get user by ID
- `PUT /v1/users/{id}` - Update user
- `DELETE /v1/users/{id}` - Delete user

## Authentication

Most endpoints require authentication using JWT Bearer tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Generating Documentation

To regenerate the Swagger documentation after making changes to the API annotations:

```bash
# Using make command
make swagger

# Or directly using swag
swag init -g cmd/server/main.go -o docs --parseDependency
```

## Example Usage

### Register a new user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "John Doe",
    "password": "password123",
    "role": "user"
  }'
```

### Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Get user profile (with authentication):
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Swagger Annotations

The API uses swaggo annotations to generate documentation. Key annotations include:

- `@Summary` - Brief description of the endpoint
- `@Description` - Detailed description
- `@Tags` - Group endpoints by tags
- `@Accept` - Request content type
- `@Produce` - Response content type
- `@Param` - Request parameters
- `@Success` - Success response
- `@Failure` - Error responses
- `@Security` - Authentication requirements
- `@Router` - Route definition

## Model Examples

All request and response models include example values to make testing easier through the Swagger UI.
