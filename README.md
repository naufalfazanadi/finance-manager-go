# Finance Manager Go

A clean architecture REST API built with Go, Fiber, and PostgreSQL for managing personal finances. This project follows Domain-Driven Design (DDD) principles and clean architecture patterns for maintainable and scalable code.

## 🚀 Features

- **Clean Architecture**: Domain-driven design with clear sepa   └─   └─   └──            │   ├── auth_manager.go         # Central auth service manager
       │   ├── jwt_service.go          # JWT token service
       │   └── password_service.go     # Password hashing service── auth/                       # Authentication infrastructure
       │   ├── auth_manager.go         # Central auth service manager
       │   ├── jwt_service.go          # JWT token service
       │   └── password_service.go     # Password hashing servicestructure/                 # Infrastructure Layer
       ├── auth/                       # Authentication infrastructure
       │   ├── auth_manager.go         # Central auth service manager
       │   ├── jwt_service.go          # JWT token service
       │   └── password_service.go     # Password hashing service
       └── database/                   # Database infrastructure
           └── postgres.go             # PostgreSQL connectiontructure/                 # Infrastructure Layer
       ├── auth/                       # Authentication infrastructure
       │   ├── auth_manager.go         # Central auth service manager
       │   ├── jwt_service.go          # JWT token service
       │   └── password_service.go     # Password hashing service
       └── database/                   # Database infrastructurestructure/                 # Infrastructure Layer
       ├── auth/                       # Authentication infrastructure
       │   ├── auth_manager.go         # Central auth service manager
       │   ├── jwt_service.go          # JWT token service
       │   └── password_service.go     # Password hashing service
       └── database/                   # Database infrastructureof concerns
- **JWT Authentication**: Secure token-based authentication with role-based access control
- **Password Security**: Bcrypt password hashing with secure storage
- **RESTful API**: Built with Fiber web framework v2.52.9
- **Database**: PostgreSQL with GORM ORM v1.30.1
- **User Management**: Complete CRUD operations with filtering and pagination
- **Middleware**: JWT authentication, CORS, error handling, and logging middleware
- **Dependency Injection**: Centralized container for managing dependencies
- **Structured Logging**: Comprehensive logging with Logrus and contextual information
- **Validation**: Request validation with go-playground/validator v10.27.0
- **Live Reload**: Development with Air (like nodemon for Go)
- **UUID Support**: Google UUID for unique identifiers
- **Environment Management**: dotenv for configuration
- **Error Handling**: Comprehensive error responses with custom AppError types
- **API Documentation**: Interactive Swagger/OpenAPI documentation with examples
- **Role-Based Access**: Admin and user roles with protected endpoints
- **Swagger UI**: Interactive API documentation accessible at `/swagger/index.html`

## 📋 Prerequisites

- Go 1.24.4 or higher
- PostgreSQL 12 or higher
- Air (for live reload development)

## 🛠️ Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/naufalfazanadi/finance-manager-go.git
   cd finance-manager-go
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Install Air for live reload**
   ```bash
   go install github.com/air-verse/air@latest
   ```

4. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` file with your configuration:
   ```properties
   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=finance_manager_db
   DB_SSLMODE=disable

   # Server
   SERVER_PORT=8080
   SERVER_HOST=localhost

   # Application
   APP_ENV=development
   LOG_LEVEL=debug

   # JWT Configuration
   JWT_SECRET=your-secret-key-change-in-production
   JWT_EXPIRES_IN=24h
   ```

5. **Setup PostgreSQL database**
   ```sql
   CREATE DATABASE finance_manager_db;
   ```

6. **Run database migrations**
   ```bash
   # Migrations will run automatically on application start
   # Check /migrations folder for SQL files
   ```

## 🏃‍♂️ Running the Application

### Development (with live reload)
```bash
# Using Air (recommended for development)
air

# Using Make command
make dev

# Using Windows batch script
dev.bat

# Using VS Code task
# Run "Air - Live Reload Server" task from VS Code
```

### Production
```bash
# Build and run
make run

# Or build separately
make build
./tmp/finance-manager

# Direct go run
go run cmd/server/main.go
```

### Available Make Commands
```bash
make dev          # Run with air (live reload)
make swagger      # Generate Swagger documentation
make build        # Build the application (includes swagger generation)
make run          # Build and run the application
make start        # Run without live reload
make clean        # Clean build artifacts
make test         # Run tests
make deps         # Download dependencies
make fmt          # Format code
make tidy         # Tidy dependencies
make help         # Show all commands
```

## 📚 API Documentation

### Swagger UI
Interactive API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

The Swagger UI provides:
- **Interactive Testing**: Test all endpoints directly from the browser
- **Request/Response Examples**: See example data for all endpoints
- **Authentication Support**: Built-in JWT token authentication
- **Model Schemas**: Detailed request and response models
- **Try It Out**: Execute real API calls with custom parameters

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check
```bash
GET /health
GET /
```

### Authentication Endpoints
```bash
# Register a new user
POST /api/v1/auth/register
Content-Type: application/json
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "password123",
  "role": "user"
}

# Login user
POST /api/v1/auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "password123"
}

# Get authenticated user profile (Protected - requires JWT token)
GET /api/v1/auth/profile
Authorization: Bearer <jwt_token>
```

### User Management Endpoints
```bash
# Create user (Public - alternative to register)
POST /api/v1/users
Content-Type: application/json
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "password123",
  "role": "user"
}

# Get all users (Protected - requires JWT token)
GET /api/v1/users?page=1&limit=10&search=john&sort_by=name&sort_type=asc
Authorization: Bearer <jwt_token>

# Get user by ID (Protected - requires JWT token)
GET /api/v1/users/{id}
Authorization: Bearer <jwt_token>

# Update user (Protected - requires JWT token)
PUT /api/v1/users/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json
{
  "name": "John Updated"
}

# Delete user (Protected - requires JWT token)
DELETE /api/v1/users/{id}
Authorization: Bearer <jwt_token>
```

### Advanced Query Features

The API supports comprehensive filtering, sorting, and pagination:

```bash
# Pagination
GET /api/v1/users?page=2&limit=10

# Search users by name or email
GET /api/v1/users?search=john

# Filter by role
GET /api/v1/users?role=admin

# Sort users (available: name, email, role, created_at, updated_at)
GET /api/v1/users?sort_by=created_at&sort_type=desc

# Date range filtering
GET /api/v1/users?created_after=2024-01-01&created_before=2024-12-31

# Combined filtering example
GET /api/v1/users?search=john&role=user&sort_by=name&sort_type=asc&page=1&limit=20
```

### Response Format

#### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Paginated Response
```json
{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

#### Authentication Response
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Error Response Format
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error information"
}
```

#### Validation Error Example
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"
}
```

## 🏗️ Project Structure (Clean Architecture)

```
finance-manager-go/
├── cmd/
│   └── server/
│       └── main.go                     # Application entry point
├── internal/
│   ├── app/                            # Application Layer
│   │   ├── container/                  # Dependency Injection Container
│   │   │   └── container.go            # Centralized dependency management
│   │   ├── handlers/                   # HTTP handlers (controllers)
│   │   │   ├── auth_handler.go         # Authentication HTTP handlers
│   │   │   └── user_handler.go         # User HTTP handlers
│   │   ├── middleware/                 # HTTP middleware
│   │   │   └── middleware.go           # JWT auth, CORS, error handling
│   │   └── routes/                     # Route definitions
│   │       ├── auth_route.go           # Authentication routes
│   │       ├── routes.go               # Main router setup
│   │       └── user_route.go           # User routes
│   ├── domain/                         # Domain Layer (Business Logic)
│   │   ├── entities/                   # Domain entities
│   │   │   └── user.go                 # User entity with roles
│   │   ├── repositories/               # Repository interfaces
│   │   │   └── user_repository.go      # User repository interface
│   │   └── usecases/                   # Business use cases
│   │       ├── auth_usecase.go         # Authentication business logic
│   │       └── user_usecase.go         # User business logic
│   ├── dto/                            # Data Transfer Objects
│   │   ├── auth_dto.go                 # Authentication DTOs
│   │   ├── common_dto.go               # Common DTOs (responses, pagination)
│   │   └── user_dto.go                 # User-specific DTOs
│   └── infrastructure/                 # Infrastructure Layer
│       ├── auth/                       # Authentication infrastructure
│       │   ├── auth_manager.go         # Central auth service manager
│       │   ├── jwt_service.go          # JWT token service
│       │   └── password_service.go     # Password hashing service
│       └── database/                   # Database infrastructure
│           └── postgres.go             # PostgreSQL connection
├── pkg/                               # Shared Packages
│   ├── config/                        # Configuration management
│   │   └── config.go                  # App configuration with JWT settings
│   ├── helpers/                       # Helper utilities
│   │   ├── errors.go                  # Custom error types and handling
│   │   └── response.go                # HTTP response utilities
│   ├── logger/                        # Logging utilities
│   │   └── logger.go                  # Structured logger with context
│   ├── types/                         # Shared types
│   │   └── database.go                # Database configuration types
│   ├── utils/                         # Common utilities
│   │   └── constant.go                # Application constants
│   └── validator/                     # Validation utilities
│       └── validator.go               # Custom validators
├── migrations/                        # Database migrations
│   ├── 003_add_password_to_users.sql  # Password field migration
│   └── 004_add_role_to_users.sql      # User roles migration
├── docs/                             # Documentation
│   └── FILTERING_EXAMPLES.md         # API filtering examples
├── scripts/                          # Build and deployment scripts
├── tmp/                              # Temporary build files
│   ├── build-errors.log              # Build error logs
│   ├── main                          # Linux binary
│   └── main.exe                      # Windows binary
├── worker/                           # Background workers (future)
├── .air.toml                         # Air configuration
├── .env.example                      # Environment template
├── AUTHENTICATION_SUMMARY.md         # Authentication implementation guide
├── DEPENDENCY_INJECTION.md           # DI architecture documentation
├── dev.bat                           # Windows development script
├── Makefile                          # Build commands
├── README.md                         # This file
├── go.mod                            # Go modules
└── go.sum                            # Go dependencies checksum
```

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`): Core business logic and entities
   - **Entities**: Business objects with business rules (User with roles)
   - **Repositories**: Interfaces for data access with filtering support
   - **Use Cases**: Application-specific business rules (Auth & User management)

2. **Application Layer** (`internal/app/`): Application services and handlers
   - **Container**: Centralized dependency injection container
   - **Handlers**: HTTP request/response handling for auth and users
   - **Middleware**: Cross-cutting concerns (JWT auth, CORS, error handling)
   - **Routes**: API endpoint definitions with authentication

3. **Infrastructure Layer** (`internal/infrastructure/`): External concerns
   - **Database**: Data persistence implementation with PostgreSQL
   - **Auth**: Authentication services (JWT, password hashing)

4. **Shared Layer** (`pkg/`): Common utilities and configurations
   - **Config**: Application configuration with JWT settings
   - **Helpers**: Error handling and HTTP response utilities
   - **Logger**: Structured logging with contextual information
   - **Types**: Common type definitions
   - **Validator**: Request validation utilities

## 🔐 Authentication Architecture

The application implements a comprehensive JWT-based authentication system:

### Authentication Features
- **User Registration**: Create accounts with email, password, and role assignment
- **User Login**: Authenticate with email/password and receive JWT tokens
- **JWT Tokens**: Secure token-based authentication with configurable expiration
- **Password Security**: Bcrypt hashing with secure storage (passwords never exposed)
- **Role-Based Access**: Admin and user roles with different permission levels
- **Protected Routes**: Middleware-based route protection using JWT validation
- **Profile Management**: Access authenticated user profile information

### Authentication Services
- **JWT Service**: Token generation, validation, and extraction from headers
- **Password Service**: Secure password hashing and verification using bcrypt
- **Auth Manager**: Centralized authentication service coordination

### Security Measures
- **Password Hashing**: All passwords hashed using bcrypt with default cost
- **Token Expiration**: JWT tokens expire after configurable time (default: 24h)
- **Secure Headers**: Proper Authorization header validation
- **Input Validation**: Request payload validation for all auth endpoints
- **Error Handling**: Consistent error responses without information leakage

## 📦 Dependency Injection Architecture

The application uses a centralized dependency injection container to manage all dependencies:

### Container Benefits
- **Single Initialization**: All dependencies initialized once in the container
- **No Duplication**: Eliminates duplicate dependency creation across modules
- **Easy Testing**: Mock entire container for comprehensive testing
- **Better Performance**: Reuse instances across the application
- **Maintainability**: Changes to dependencies happen in one place

### Dependency Flow
```
Container → Repositories → Infrastructure → Use Cases → Handlers → Routes
```

### Available Dependencies
- **Repositories**: User repository with database operations
- **Infrastructure**: Auth manager, JWT service, password service
- **Use Cases**: Authentication and user business logic
- **Handlers**: HTTP request handlers for auth and user endpoints
- **Middleware**: JWT authentication and authorization middleware

## 📖 Additional Documentation

For more detailed information, check out these documentation files:

- **[AUTHENTICATION_SUMMARY.md](AUTHENTICATION_SUMMARY.md)**: Complete guide to the JWT authentication implementation
- **[DEPENDENCY_INJECTION.md](DEPENDENCY_INJECTION.md)**: Architecture overview of the centralized dependency injection system
- **[docs/FILTERING_EXAMPLES.md](docs/FILTERING_EXAMPLES.md)**: Comprehensive examples of API filtering, sorting, and pagination

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests for specific package
go test ./internal/domain/usecases/...

# Run tests with verbose output
go test -v ./...
```

## 📦 Dependencies

### Core Dependencies
- **Fiber v2.52.9**: Fast HTTP web framework
- **GORM v1.30.1**: Go ORM library
- **PostgreSQL Driver v1.6.0**: GORM PostgreSQL driver
- **Logrus v1.9.3**: Structured logging
- **Validator v10.27.0**: Request validation
- **JWT v5.3.0**: JSON Web Token implementation
- **Bcrypt**: Password hashing for security
- **UUID v1.6.0**: UUID generation
- **Godotenv v1.5.1**: Environment variable loading

### Development Dependencies
- **Air**: Live reload for development

## � Configuration

Environment variables are managed through `.env` file:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | Database host | `localhost` | Yes |
| `DB_PORT` | Database port | `5432` | Yes |
| `DB_USER` | Database username | `postgres` | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | `finance_manager_db` | Yes |
| `DB_SSLMODE` | SSL mode | `disable` | Yes |
| `SERVER_PORT` | Server port | `8080` | Yes |
| `SERVER_HOST` | Server host | `localhost` | Yes |
| `APP_ENV` | Environment (development/production) | `development` | Yes |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `debug` | Yes |
| `JWT_SECRET` | JWT signing secret | `your-secret-key-change-in-production` | Yes |
| `JWT_EXPIRES_IN` | JWT token expiration | `24h` | Yes |

### Environment-Specific Settings

#### Development
```env
APP_ENV=development
LOG_LEVEL=debug
DB_SSLMODE=disable
JWT_SECRET=dev-secret-key-not-for-production
```

#### Production
```env
APP_ENV=production
LOG_LEVEL=info
DB_SSLMODE=require
JWT_SECRET=your-super-secure-random-secret-key
```

## 🚀 Deployment

### Build for Production
```bash
# Build optimized binary
go build -ldflags="-w -s" -o finance-manager cmd/server/main.go

# Or use Makefile
make build
```

### Production Security Checklist
- [ ] Change JWT_SECRET to a strong, random secret key
- [ ] Set APP_ENV to "production"
- [ ] Use strong database passwords
- [ ] Enable SSL/TLS for database connections (set DB_SSLMODE=require)
- [ ] Set up proper firewall rules
- [ ] Use environment variables for all secrets
- [ ] Enable request logging and monitoring
- [ ] Set up database backups

### Docker (Coming Soon)
```bash
# Build and run with Docker Compose
docker-compose up --build

# Run in detached mode
docker-compose up -d
```

## 🔍 Development Tools

- **Air**: Live reload during development
- **VS Code Tasks**: Pre-configured development tasks
- **Makefile**: Common development commands
- **Windows Batch Scripts**: Windows-specific development commands

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Follow the clean architecture principles
4. Write tests for your code
5. Ensure all tests pass (`make test`)
6. Format your code (`make fmt`)
7. Commit your changes (`git commit -m 'Add some amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Code Style Guidelines
- Follow Go naming conventions
- Use meaningful variable and function names
- Keep functions small and focused
- Write comprehensive tests
- Document public APIs
- Follow clean architecture principles

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Fiber](https://github.com/gofiber/fiber) - Express-inspired web framework
- [GORM](https://github.com/go-gorm/gorm) - The fantastic ORM library for Golang
- [Logrus](https://github.com/sirupsen/logrus) - Structured, pluggable logging
- [Air](https://github.com/air-verse/air) - Live reload utility for Go apps
- [Go Playground Validator](https://github.com/go-playground/validator) - Go struct and field validation

## 📧 Contact

For questions or support, please open an issue on GitHub or contact the maintainer.

---

Built with ❤️ using Go and Clean Architecture principles