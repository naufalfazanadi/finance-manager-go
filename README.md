# Finance Manager Go

A comprehensive clean architecture REST API built with Go, Fiber, and PostgreSQL for managing personal finances. This project follows Domain-Driven Design (DDD) principles and clean architecture patterns for maintainable and scalable code.

## 🚀 Features

### 💰 Finance Management
- **Transaction Management**: Complete CRUD operations for income and expense transactions
- **Wallet Management**: Complete CRUD operations for personal wallets with different types and categories
- **Multi-Currency Support**: Handle different currencies (IDR, USD, EUR, etc.)
- **Balance Tracking**: Track wallet balances with decimal precision and automatic updates
- **Transaction Categories**: Categorize transactions for better organization
- **Transaction Types**: Support for income a## ⚙️ Configuration

Environment variables are managed through `.env` file. Below are all available configuration options:

### Database Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_HOST` | Database host | `localhost` | Yes |
| `DB_PORT` | Database port | `5432` | Yes |
| `DB_USER` | Database username | `postgres` | Yes |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_NAME` | Database name | `clean_api_db` | Yes |
| `DB_SSLMODE` | SSL mode | `disable` | Yes |

### Database Connection Pool Settings
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_MAX_OPEN_CONNS` | Maximum open connections | `25` | No |
| `DB_MAX_IDLE_CONNS` | Maximum idle connections | `5` | No |
| `DB_CONN_MAX_LIFETIME` | Connection lifetime (minutes) | `30` | No |
| `DB_CONN_MAX_IDLE_TIME` | Idle timeout (minutes) | `5` | No |
| `DB_MAX_RETRIES` | Retry attempts | `3` | No |
| `DB_RETRY_DELAY` | Delay between retries (seconds) | `5` | No |
| `DB_CONNECT_TIMEOUT` | Initial connection timeout (seconds) | `10` | No |

### Server Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SERVER_PORT` | Server port | `8080` | Yes |
| `SERVER_HOST` | Server host | `localhost` | Yes |

### Application Settings
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `APP_ENV` | Environment (development/production) | `development` | Yes |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | `debug` | Yes |

### JWT Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `JWT_SECRET` | JWT signing secret | `your-secret-key-change-in-production` | Yes |
| `JWT_EXPIRES_IN` | JWT token expiration | `24h` | Yes |

### CORS Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `CORS_ALLOW_ORIGINS` | Allowed origins | `*` | No |
| `CORS_ALLOW_METHODS` | Allowed methods | `GET,POST,PUT,DELETE,OPTIONS` | No |
| `CORS_ALLOW_HEADERS` | Allowed headers | `Origin,Content-Type,Accept,Authorization` | No |

### MinIO Configuration
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `MINIO_ENDPOINT` | MinIO server endpoint | `localhost:9000` | No |
| `MINIO_ACCESS_KEY` | MinIO access key | `minioadmin` | No |
| `MINIO_SECRET_KEY` | MinIO secret key | `minioadmin` | No |
| `MINIO_USE_SSL` | Use SSL for MinIO | `false` | No |
### Email (SMTP) Configuration
SMTP settings are now managed via the application's config system (`pkg/config/config.go`).
You can set these via environment variables in your `.env` file, and they will be loaded into the config struct at startup:

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SMTP_HOST` | SMTP server host | `smtp.gmail.com` | Yes |
| `SMTP_PORT` | SMTP server port | `587` | Yes |
| `SMTP_USERNAME` | SMTP username | - | Yes |
| `SMTP_PASSWORD` | SMTP password or app password | - | Yes |
| `SMTP_FROM_EMAIL` | From email address | - | Yes |
| `SMTP_FROM_NAME` | From name | `Finance Manager` | No |
| `FRONTEND_URL` | Frontend URL for reset links | `http://localhost:3000` | No |

**Note:**
- The application uses the `SMTPConfig` struct in `pkg/config/config.go` to manage all SMTP settings.
- The mail sending logic (`pkg/mail/mail.go`) now always reads SMTP config from the config package, not directly from environment variables.
- Update your `.env` file as before; the config loader will handle the rest.
| `MINIO_PRIVATE_BUCKET` | Private bucket name | `private` | No |
| `MINIO_PUBLIC_BUCKET` | Public bucket name | `public` | No |
| `MINIO_DIRECTORY` | Directory prefix | `` | No |

### Email Configuration (For Password Reset)
| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `SMTP_HOST` | SMTP server host | `smtp.gmail.com` | Yes |
| `SMTP_PORT` | SMTP server port | `587` | Yes |
| `SMTP_USERNAME` | SMTP username | - | Yes |
| `SMTP_PASSWORD` | SMTP password or app password | - | Yes |
| `SMTP_FROM_EMAIL` | From email address | - | Yes |
| `SMTP_FROM_NAME` | From name | `Finance Manager` | No |
| `FRONTEND_URL` | Frontend URL for reset links | `http://localhost:3000` | No |tions with proper wallet impact calculation
- **Soft Delete**: Recoverable deletion with restore functionality for both transactions and wallets

### 🔐 Authentication & Security
- **JWT Authentication**: Secure token-based authentication with configurable expiration
- **Strong Password Validation**: Password strength requirements with custom validation
- **Password Security**: Bcrypt password hashing with secure storage
- **Role-Based Access**: Admin and user roles with protected endpoints
- **PII Encryption**: Advanced encryption for user emails and sensitive information
- **Rate Limiting**: Built-in rate limiting middleware for API protection

### 👥 User Management
- **Complete User CRUD**: Create, read, update, delete operations with validation
- **Profile Photos**: Upload and manage user profile photos with MinIO integration
- **User Filtering**: Advanced search, sorting, and pagination capabilities
- **Soft Delete Support**: Recoverable user deletion with restore functionality
- **Birth Date Management**: Encrypted birth date storage with age calculation
- **DataDog Integration**: Built-in monitoring and tracing capabilities

### 📁 File Management
- **MinIO Integration**: Secure file storage with public/private bucket support
- **Profile Photo Upload**: Support for JPEG, PNG formats with size validation
- **File Validation**: Comprehensive file type and size validation
- **Automatic Cleanup**: Failed upload rollback and old file cleanup
- **Upload Helper**: Centralized upload utilities with validation configs

### 🏗️ Technical Features
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **RESTful API**: Built with Fiber web framework v2.52.9
- **Database**: PostgreSQL with GORM ORM v1.30.1 and optimized connection pooling
- **Middleware**: JWT authentication, CORS, error handling, rate limiting, and logging middleware
- **Dependency Injection**: Centralized container for managing all application dependencies
- **Structured Logging**: Comprehensive logging with Logrus and contextual information
- **Advanced Validation**: Request validation with custom rules and file validation
- **Live Reload**: Development with Air (like nodemon for Go)
- **UUID Support**: Google UUID for unique identifiers throughout the system
- **Environment Management**: Comprehensive configuration with dotenv
- **Error Handling**: Comprehensive error responses with custom AppError types
- **API Documentation**: Interactive Swagger/OpenAPI documentation with examples
- **Generic Pagination**: Type-safe pagination system with metadata
- **Swagger UI**: Interactive API documentation accessible at `/swagger/index.html`
- **Query Parser**: Advanced query parsing for filtering and sorting
- **Health Checks**: Built-in health check endpoints

## 📋 Prerequisites

- Go 1.24.4 or higher
- PostgreSQL 12 or higher
- MinIO (for file storage)
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
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=finance_manager_db
   DB_SSLMODE=disable

   # Database Connection Pool Settings
   DB_MAX_OPEN_CONNS=25
   DB_MAX_IDLE_CONNS=5
   DB_CONN_MAX_LIFETIME=30
   DB_CONN_MAX_IDLE_TIME=5
   DB_MAX_RETRIES=3
   DB_RETRY_DELAY=5
   DB_CONNECT_TIMEOUT=10

   # Server Configuration
   SERVER_PORT=8080
   SERVER_HOST=localhost

   # Application Settings
   APP_ENV=development
   LOG_LEVEL=debug

   # JWT Configuration
   JWT_SECRET=your-secret-key-change-in-production
   JWT_EXPIRES_IN=24h

   # CORS Configuration
   CORS_ALLOW_ORIGINS=*
   CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
   CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization

   # MinIO Configuration
   MINIO_ENDPOINT=localhost:9000
   MINIO_ACCESS_KEY=minioadmin
   MINIO_SECRET_KEY=minioadmin
   MINIO_USE_SSL=false
   MINIO_PRIVATE_BUCKET=private
   MINIO_PUBLIC_BUCKET=public
   MINIO_DIRECTORY=finance-manager
   ```

5. **Setup PostgreSQL database**
   ```sql
   CREATE DATABASE finance_manager_db;
   ```

6. **Setup MinIO (for file storage)**
   ```bash
   # Using Docker
   docker run -d -p 9000:9000 -p 9001:9001 \
     --name minio \
     -e "MINIO_ROOT_USER=minioadmin" \
     -e "MINIO_ROOT_PASSWORD=minioadmin" \
     minio/minio server /data --console-address ":9001"
   ```

7. **Run database migrations**
   ```bash
   # Database schema will be auto-migrated on application start
   # Check internal/domain/entities/ for entity definitions
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
make install-air  # Install air for development
make fmt          # Format code
make lint         # Lint code (requires golangci-lint)
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
  "password": "Password123!",
  "birth_date": "1990-01-15"
}

# Login user
POST /api/v1/auth/login
Content-Type: application/json
{
  "email": "user@example.com",
  "password": "Password123!"
}

# Get authenticated user profile (Protected - requires JWT token)
GET /api/v1/auth/profile
Authorization: Bearer <jwt_token>
```

### User Management Endpoints
```bash
# Create user (Protected - admin only)
POST /api/v1/users
Content-Type: multipart/form-data
# OR Content-Type: application/json
{
  "email": "user@example.com",
  "name": "John Doe",
  "password": "Password123!",
  "birth_date": "1990-01-15"
}
# Optional: profile_photo_file (multipart file)

# Get all users (Protected - requires JWT token)
GET /api/v1/users?page=1&limit=10&search=john&sort_by=name&sort_type=asc
Authorization: Bearer <jwt_token>

# Get user by ID (Protected - requires JWT token)
GET /api/v1/users/{id}
Authorization: Bearer <jwt_token>

# Update user (Protected - user can update own profile, admin can update any)
PUT /api/v1/users/{id}
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data
# OR Content-Type: application/json
{
  "name": "John Updated",
  "birth_date": "1990-01-15"
}
# Optional: profile_photo_file (multipart file)

# Delete user (Protected - soft delete)
DELETE /api/v1/users/{id}
Authorization: Bearer <jwt_token>

# Restore user (Protected - admin only)
PUT /api/v1/users/{id}/restore
Authorization: Bearer <jwt_token>

# Hard delete user (Protected - admin only)
DELETE /api/v1/users/{id}/hard
Authorization: Bearer <jwt_token>

# Get users with deleted (Protected - admin only)
GET /api/v1/users/with-deleted?page=1&limit=10
Authorization: Bearer <jwt_token>

# Get only deleted users (Protected - admin only)
GET /api/v1/users/deleted?page=1&limit=10
Authorization: Bearer <jwt_token>
```

### Wallet Management Endpoints
```bash
# Create wallet (Protected - user can create own wallet, admin can create for any user)
POST /api/v1/wallets
Authorization: Bearer <jwt_token>
Content-Type: application/json
{
  "name": "My Savings",
  "type": "savings",
  "category": "personal",
  "balance": 1000.50,
  "currency": "IDR",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}

# Get all wallets (Protected - user sees own wallets, admin sees all)
GET /api/v1/wallets?page=1&limit=10&search=savings&sort_by=created_at&sort_type=desc
Authorization: Bearer <jwt_token>

# Get wallet by ID (Protected - user can access own wallets, admin can access all)
GET /api/v1/wallets/{id}
Authorization: Bearer <jwt_token>

# Update wallet (Protected - user can update own wallets, admin can update any)
PUT /api/v1/wallets/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json
{
  "name": "Updated Savings",
  "balance": 2000.75
}

# Delete wallet (Protected - soft delete)
DELETE /api/v1/wallets/{id}
Authorization: Bearer <jwt_token>

# Restore wallet (Protected - user can restore own wallets, admin can restore any)
PUT /api/v1/wallets/{id}/restore
Authorization: Bearer <jwt_token>
```

### Transaction Management Endpoints
```bash
# Create transaction (Protected - user can create own transactions, admin can create for any user)
POST /api/v1/transactions
Authorization: Bearer <jwt_token>
Content-Type: application/json
{
  "name": "Grocery Shopping",
  "cost": 150.75,
  "type": "expense",
  "note": "Weekly grocery shopping",
  "t_category": "food",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "wallet_id": "123e4567-e89b-12d3-a456-426614174001"
}

# Get all transactions (Protected - user sees own transactions, admin sees all)
GET /api/v1/transactions?page=1&limit=10&search=grocery&sort_by=created_at&sort_type=desc
Authorization: Bearer <jwt_token>

# Get transaction by ID (Protected - user can access own transactions, admin can access all)
GET /api/v1/transactions/{id}
Authorization: Bearer <jwt_token>

# Update transaction (Protected - user can update own transactions, admin can update any)
PUT /api/v1/transactions/{id}
Authorization: Bearer <jwt_token>
Content-Type: application/json
{
  "name": "Updated Grocery Shopping",
  "cost": 175.25,
  "note": "Weekly grocery shopping with extra items"
}

# Delete transaction (Protected - soft delete)
DELETE /api/v1/transactions/{id}
Authorization: Bearer <jwt_token>

# Restore transaction (Protected - user can restore own transactions, admin can restore any)
PUT /api/v1/transactions/{id}/restore
Authorization: Bearer <jwt_token>

# Get transactions by wallet (Protected)
GET /api/v1/transactions/wallet/{wallet_id}?page=1&limit=10
Authorization: Bearer <jwt_token>

# Get transactions by type (Protected)
GET /api/v1/transactions?type=income&page=1&limit=10
GET /api/v1/transactions?type=expense&page=1&limit=10
Authorization: Bearer <jwt_token>
```

### Advanced Query Features

The API supports comprehensive filtering, sorting, and pagination across all endpoints:

```bash
# Pagination (all endpoints)
GET /api/v1/users?page=2&limit=10
GET /api/v1/wallets?page=1&limit=5
GET /api/v1/transactions?page=1&limit=20

# Search functionality
GET /api/v1/users?search=john                # Search users by name or email
GET /api/v1/wallets?search=savings           # Search wallets by name, type, or category
GET /api/v1/transactions?search=grocery      # Search transactions by name or note

# Filter by specific fields
GET /api/v1/users?role=admin
GET /api/v1/wallets?category=personal&type=savings
GET /api/v1/transactions?type=income&t_category=salary

# Sort options (available fields vary by endpoint)
# Users: name, email, created_at, updated_at
GET /api/v1/users?sort_by=created_at&sort_type=desc

# Wallets: name, type, category, balance, created_at, updated_at
GET /api/v1/wallets?sort_by=balance&sort_type=asc

# Transactions: name, cost, type, created_at, updated_at
GET /api/v1/transactions?sort_by=cost&sort_type=desc

# Date range filtering
GET /api/v1/users?created_after=2024-01-01&created_before=2024-12-31
GET /api/v1/wallets?created_after=2024-01-01&created_before=2024-12-31
GET /api/v1/transactions?created_after=2024-01-01&created_before=2024-12-31

# Combined filtering examples
GET /api/v1/wallets?search=savings&category=personal&sort_by=balance&sort_type=desc&page=1&limit=10
GET /api/v1/transactions?type=expense&t_category=food&sort_by=cost&sort_type=desc&page=1&limit=10

# Soft delete filtering (admin only)
GET /api/v1/users/with-deleted?page=1&limit=10    # Include soft deleted records
GET /api/v1/users/deleted?page=1&limit=10         # Only soft deleted records
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
    "birth_date": "1990-01-15T00:00:00Z",
    "age": 34,
    "profile_photo": "https://minio.example.com/public/profile-photo/2024/01/profile_photo_1641024000.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Transaction Response
```json
{
  "success": true,
  "message": "Transaction retrieved successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "Grocery Shopping",
    "cost": 150.75,
    "type": "expense",
    "note": "Weekly grocery shopping",
    "t_category": "food",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "wallet_id": "123e4567-e89b-12d3-a456-426614174001",
    "is_deleted": false,
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "John Doe",
      "email": "user@example.com"
    },
    "wallet": {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "name": "My Checking Account",
      "type": "checking",
      "currency": "IDR"
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Wallet Response
```json
{
  "success": true,
  "message": "Wallet retrieved successfully",
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "My Savings",
    "type": "savings",
    "category": "personal",
    "balance": 1000.50,
    "currency": "IDR",
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "is_deleted": false,
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "John Doe",
      "email": "user@example.com"
    },
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
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user",
    "birth_date": "1990-01-15T00:00:00Z",
    "age": 34,
    "profile_photo": "https://minio.example.com/public/profile-photo/2024/01/profile_photo_1641024000.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z",
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
  "message": "Form validation failed",
  "error": "field 'password' must contain at least 1 uppercase letter, 1 number, and 1 special character"
}
```

#### File Upload Error Example
```json
{
  "success": false,
  "message": "File validation failed",
  "error": "File size must not exceed 2097152 bytes"
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
│   │   │   └── service_container.go    # Centralized dependency management
│   │   ├── handlers/                   # HTTP handlers (controllers)
│   │   │   ├── auth_handler.go         # Authentication HTTP handlers
│   │   │   ├── health_handler.go       # Health check handlers
│   │   │   ├── transaction_handler.go  # Transaction HTTP handlers
│   │   │   ├── user_handler.go         # User HTTP handlers
│   │   │   └── wallet_handler.go       # Wallet HTTP handlers
│   │   ├── middleware/                 # HTTP middleware
│   │   │   ├── middleware.go           # JWT auth, CORS, error handling
│   │   │   └── rate_limiter.go         # Rate limiting middleware
│   │   └── routes/                     # Route definitions
│   │       ├── auth_route.go           # Authentication routes
│   │       ├── routes.go               # Main router setup
│   │       ├── transaction_route.go    # Transaction routes
│   │       ├── user_route.go           # User routes
│   │       └── wallet_route.go         # Wallet routes
│   ├── domain/                         # Domain Layer (Business Logic)
│   │   ├── entities/                   # Domain entities
│   │   │   ├── transaction.go          # Transaction entity with types and calculations
│   │   │   ├── user.go                 # User entity with roles and encryption
│   │   │   └── wallet.go               # Wallet entity with soft delete
│   │   ├── repositories/               # Repository interfaces
│   │   │   ├── transaction_repository.go # Transaction repository interface
│   │   │   ├── user_repository.go      # User repository interface
│   │   │   └── wallet_repository.go    # Wallet repository interface
│   │   └── usecases/                   # Business use cases
│   │       ├── auth_usecase.go         # Authentication business logic
│   │       ├── transaction_usecase.go  # Transaction business logic
│   │       ├── user_usecase.go         # User business logic
│   │       └── wallet_usecase.go       # Wallet business logic
│   ├── dto/                            # Data Transfer Objects
│   │   ├── auth_dto.go                 # Authentication DTOs
│   │   ├── common_dto.go               # Common DTOs (responses, pagination)
│   │   ├── transaction_dto.go          # Transaction-specific DTOs
│   │   ├── user_dto.go                 # User-specific DTOs
│   │   └── wallet_dto.go               # Wallet-specific DTOs
│   └── infrastructure/                 # Infrastructure Layer
│       ├── auth/                       # Authentication infrastructure
│       │   ├── jwt.go                  # JWT token service
│       │   └── password.go             # Password hashing with validation
│       ├── cache/                      # Cache infrastructure (future)
│       └── database/                   # Database infrastructure
│           └── postgres.go             # PostgreSQL connection with pooling
├── pkg/                               # Shared Packages
│   ├── config/                        # Configuration management
│   │   └── config.go                  # App configuration with all settings
│   ├── encryption/                    # PII encryption utilities
│   │   └── pii_encryption.go         # Email and sensitive data encryption
│   ├── helpers/                       # Helper utilities
│   │   ├── errors.go                  # Custom error types and handling
│   │   ├── query_parser.go            # Query parameter parsing utilities
│   │   └── response.go                # HTTP response utilities
│   ├── logger/                        # Logging utilities
│   │   └── logger.go                  # Structured logger with context
│   ├── minio/                         # MinIO file storage
│   │   ├── client.go                  # MinIO client initialization
│   │   ├── download.go                # File download operations
│   │   ├── helper.go                  # Upload/delete helper functions
│   │   └── upload.go                  # File upload operations
│   ├── upload/                        # File validation
│   │   └── validation_configs.go      # File validation configurations
│   ├── utils/                         # Common utilities
│   │   └── constant.go                # Application constants and messages
│   └── validator/                     # Validation utilities
│       └── validator.go               # Custom validators with file support
├── migrations/                        # Database migrations (currently empty - auto-migration used)
├── docs/                             # Documentation
│   ├── docs.go                       # Swagger docs generation
│   ├── swagger.json                  # Generated Swagger JSON
│   ├── swagger.yaml                  # Generated Swagger YAML
│   └── sequence/                     # Sequence diagrams (future)
├── examples/                         # API usage examples (future)
├── scripts/                          # Build and deployment scripts
├── storage/                          # Local file storage
│   ├── private/                      # Private files
│   └── public/                       # Public files
├── tmp/                              # Temporary build files
│   ├── build-errors.log              # Build error logs
│   ├── finance-manager               # Linux binary
│   ├── main                          # Linux binary (alternative)
│   ├── main.exe                      # Windows binary
│   ├── server.exe                    # Windows server binary
│   └── test-build                    # Test build artifacts
├── worker/                           # Background workers (future)
├── .air.toml                         # Air configuration
├── .env.example                      # Environment template
├── dev.bat                           # Windows development script
├── Dockerfile                        # Docker configuration
├── Makefile                          # Build commands
├── project-structure.txt             # Project structure documentation
├── README.md                         # This file
├── go.mod                            # Go modules
└── go.sum                            # Go dependencies checksum
```

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`): Core business logic and entities
   - **Entities**: Business objects with business rules (User, Wallet, Transaction)
   - **Repositories**: Interfaces for data access with advanced filtering support
   - **Use Cases**: Application-specific business rules (Auth, User, Wallet, Transaction management)

2. **Application Layer** (`internal/app/`): Application services and handlers
   - **Container**: Centralized dependency injection container
   - **Handlers**: HTTP request/response handling for all entities
   - **Middleware**: Cross-cutting concerns (JWT auth, CORS, rate limiting, error handling)
   - **Routes**: API endpoint definitions with authentication

3. **Infrastructure Layer** (`internal/infrastructure/`): External concerns
   - **Database**: Data persistence implementation with PostgreSQL and connection pooling
   - **Auth**: Authentication services (JWT, password hashing)

4. **Shared Layer** (`pkg/`): Common utilities and configurations
   - **Config**: Comprehensive application configuration with environment management
   - **Helpers**: Error handling, query parsing, and HTTP response utilities
   - **Logger**: Structured logging with contextual information
   - **MinIO**: File storage service with upload/download capabilities
   - **Validator**: Request validation utilities with file support
   - **Encryption**: PII data encryption for sensitive information

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

### Security Measures
- **Password Hashing**: All passwords hashed using bcrypt with default cost
- **Token Expiration**: JWT tokens expire after configurable time (default: 24h)
- **Secure Headers**: Proper Authorization header validation
- **Input Validation**: Request payload validation for all auth endpoints
- **Error Handling**: Consistent error responses without information leakage
- **Rate Limiting**: Built-in rate limiting to prevent abuse

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
- **Repositories**: User, Wallet, and Transaction repositories with database operations
- **Infrastructure**: Database connections, MinIO client, authentication services
- **Use Cases**: Authentication, user, wallet, and transaction business logic
- **Handlers**: HTTP request handlers for all endpoints
- **Middleware**: JWT authentication, authorization, and rate limiting middleware

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
- **MinIO Client v7.0.95**: Object storage client for file management
- **Swaggo**: Swagger documentation generation
- **DataDog Tracing v1.74.3**: Application performance monitoring and tracing

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
MINIO_USE_SSL=true
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
- [ ] Enable SSL for MinIO (set MINIO_USE_SSL=true)
- [ ] Configure proper MinIO access credentials
- [ ] Set up proper firewall rules
- [ ] Use environment variables for all secrets
- [ ] Configure database connection pool settings for production load
- [ ] Enable request logging and monitoring
- [ ] Set up database backups
- [ ] Configure DataDog tracing for production monitoring

### Docker (Coming Soon)
```bash
# Build and run with Docker Compose
docker-compose up --build

# Run in detached mode
docker-compose up -d
```

## 🔍 Development Tools

- **Air**: Live reload during development with `.air.toml` configuration
- **VS Code Tasks**: Pre-configured development tasks including live reload server
- **Makefile**: Common development commands with help documentation
- **Windows Batch Scripts**: Windows-specific development commands (`dev.bat`)
- **Swagger**: Interactive API documentation generation
- **DataDog**: Application performance monitoring and tracing

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