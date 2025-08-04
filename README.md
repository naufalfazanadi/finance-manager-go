# Finance Manager Go

A clean architecture REST API built with Go, Fiber, and PostgreSQL for managing personal finances. This project follows Domain-Driven Design (DDD) principles and clean architecture patterns for maintainable and scalable code.

## 🚀 Features

- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **RESTful API**: Built with Fiber web framework v2.52.9
- **Database**: PostgreSQL with GORM ORM v1.30.1
- **Logging**: Structured logging with Logrus
- **Validation**: Request validation with go-playground/validator v10.27.0
- **Live Reload**: Development with Air (like nodemon for Go)
- **UUID Support**: Google UUID for unique identifiers
- **Environment Management**: dotenv for configuration
- **Error Handling**: Comprehensive error responses with custom DTOs
- **Modular Design**: Well-organized domain, usecase, and infrastructure layers

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
make build        # Build the application
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

### Base URL
```
http://localhost:8080/api/v1
```

### Health Check
```bash
GET /health
GET /
```

### User Management Endpoints
```bash
# Create user
POST /api/v1/users
Content-Type: application/json
{
  "email": "user@example.com",
  "name": "John Doe",
  "age": 30
}

# Get all users (with pagination and filtering)
GET /api/v1/users?page=1&limit=10&search=john&sort_by=name&sort_dir=asc

# Get user by ID
GET /api/v1/users/{id}

# Update user
PUT /api/v1/users/{id}
Content-Type: application/json
{
  "name": "John Updated",
  "age": 31
}

# Delete user
DELETE /api/v1/users/{id}
```

### Response Format
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {...},
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### Error Response Format
```json
{
  "success": false,
  "message": "Error description",
  "errors": [...],
  "data": null
}
```

## 🏗️ Project Structure (Clean Architecture)

```
finance-manager-go/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── app/                        # Application Layer
│   │   ├── handlers/               # HTTP handlers (controllers)
│   │   │   └── user_handler.go     # User HTTP handlers
│   │   ├── middleware/             # HTTP middleware
│   │   │   └── middleware.go       # Custom middleware
│   │   └── routes/                 # Route definitions
│   │       ├── routes.go           # Main router setup
│   │       └── user_route.go       # User routes
│   ├── domain/                     # Domain Layer (Business Logic)
│   │   ├── entities/               # Domain entities
│   │   │   └── user.go             # User entity
│   │   ├── repositories/           # Repository interfaces
│   │   │   └── user_repository.go  # User repository interface
│   │   └── usecases/               # Business use cases
│   │       └── user_usecase.go     # User business logic
│   ├── dto/                        # Data Transfer Objects
│   │   ├── common_dto.go           # Common DTOs (responses, pagination)
│   │   └── user_dto.go             # User-specific DTOs
│   └── infrastructure/             # Infrastructure Layer
│       └── database/               # Database infrastructure
│           └── postgres.go         # PostgreSQL connection
├── pkg/                           # Shared Packages
│   ├── config/                    # Configuration management
│   │   └── config.go              # App configuration
│   ├── helpers/                   # Helper utilities
│   ├── logger/                    # Logging utilities
│   │   └── logger.go              # Logger setup
│   ├── types/                     # Shared types
│   │   └── database.go            # Database types
│   ├── utils/                     # Common utilities
│   │   ├── constant.go            # Application constants
│   │   └── response.go            # Response utilities
│   └── validator/                 # Validation utilities
│       └── validator.go           # Custom validators
├── migrations/                    # Database migrations
├── scripts/                       # Build and deployment scripts
├── tmp/                          # Temporary build files
│   ├── build-errors.log          # Build error logs
│   ├── main                      # Linux binary
│   └── main.exe                  # Windows binary
├── docs/                         # Documentation
├── worker/                       # Background workers (future)
├── .air.toml                     # Air configuration
├── .env.example                  # Environment template
├── dev.bat                       # Windows development script
├── Makefile                      # Build commands
├── README.md                     # This file
├── go.mod                        # Go modules
└── go.sum                        # Go dependencies checksum
```

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`): Core business logic and entities
   - Entities: Business objects with business rules
   - Repositories: Interfaces for data access
   - Use Cases: Application-specific business rules

2. **Application Layer** (`internal/app/`): Application services and handlers
   - Handlers: HTTP request/response handling
   - Middleware: Cross-cutting concerns like CORS, error handling
   - Routes: API endpoint definitions

3. **Infrastructure Layer** (`internal/infrastructure/`): External concerns
   - Database: Data persistence implementation

4. **Shared Layer** (`pkg/`): Common utilities and configurations
   - Config: Application configuration
   - Utils: Shared utilities and helpers
   - Types: Common type definitions

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

## 🚀 Deployment

### Build for Production
```bash
# Build optimized binary
go build -ldflags="-w -s" -o finance-manager cmd/server/main.go

# Or use Makefile
make build
```

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