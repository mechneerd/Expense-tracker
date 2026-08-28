# Expense Tracker API - Project Structure

## Root Directory
```
expense-tracker-api/
├── go.mod                      # Module dependencies
├── go.sum                      # Module checksums
├── sqlc.yaml                   # SQL code generation config
├── main.go                     # Entry point & API routing
├── config/                     # Environment configuration
│   └── config.go               # Config loading from env vars
├── migrations/                 # SQL schema migrations
│   ├── 001_initial_schema.up.sql
│   └── 001_initial_schema.down.sql
├── internal/                   # Domain modules (internal packages)
│   ├── users/                  # User management module
│   │   ├── model/              # User model
│   │   ├── repository/         # User repository
│   │   └── handler/            # User HTTP handlers
│   ├── families/               # Family management module
│   │   ├── model/              # Family model
│   │   ├── repository/         # Family repository
│   │   └── handler/            # Family HTTP handlers
│   ├── transactions/           # Transaction module
│   │   ├── model/              # Transaction model
│   │   ├── repository/         # Transaction repository
│   │   └── handler/            # Transaction HTTP handlers
│   ├── dashboard/              # Family head dashboard module
│   │   └── dashboard.go        # Dashboard data handler
│   ├── categories/             # Expense categories module
│   │   ├── model/              # Category model
│   │   └── repository/         # Category repository
│   ├── paymentmethods/         # Payment methods module
│   │   ├── model/              # Payment method model
│   │   └── repository/         # Payment method repository
│   └── upiapps/                # UPI apps module
│       └── model/              # UPI app model
├── pkg/                        # Public packages (shared across modules)
│   ├── handler/                # Shared handler utilities
│   ├── response/               # Response formatting (HTTPResponse)
│   ├── log/                    # Logging utilities
│   ├── users/                  # User-related shared code
│   ├── transactions/           # Transaction-related shared code
│   ├── families/               # Family-related shared code
│   ├── categories/             # Category-related shared code
│   ├── paymentmethods/         # Payment method shared code
│   └── upiapps/                # UPI app shared code
├── db/                         # Database utilities (pgx pool)
└── queries/                    # Generated SQL queries (sqlc)
```

## Key Files Overview

### main.go
- Entry point for the Go API server
- Configures router (chi) with middleware (Logger, RequestID, Recoverer)
- Routes all API v1 endpoints
- Starts HTTP server on port :8080

### config/config.go
- Loads configuration from environment variables
- Database connection (pgx)
- JWT secret key, Redis config, OTP lifetime
- `MustGetDB()` helper for easy DB access

### sqlc.yaml
- Code generation configuration
- Generates Go models from SQL queries
- Output: `models/%s.sqlc.go` and `queries/%s.sqlc.go`

### migrations/001_initial_schema.up.sql
- Initial database schema with all tables:
  - users, roles, families, family_members
  - family_invitations, transaction_categories
  - payment_methods, upi_apps
  - transactions, otp_verifications, audit_logs
- Indexes for performance
- Enum types (transaction_type)

### internal/[module]/handler/
- HTTP request handlers for each domain
- Contains business logic, API responses
- Example: userHandler, familyHandler, transactionHandler

### internal/[module]/repository/
- Data access layer
- SQL queries (or interfaces for sqlc-generated)
- CRUD operations for each entity

### internal/[module]/model/
- Data structures (Go structs)
- JSON serialization annotations
- Business logic methods (where applicable)

### pkg/response/Response.go
- Standardized HTTP response format
- `HTTPResponse` struct with Success, Message, Data, Error fields
- `JSON()` helper function for consistent responses

### pkg/log/log.go
- Structured logging setup
- Request ID correlation
- Request-from-context logger

### pkg/[module]/repository/
- Shared repository interfaces and implementations
- Can be extended by internal modules

## Dependencies (go.mod)
```
module expense-tracker-api

go 1.26.3

require (
    github.com/go-chi/chi/v5 v5.3.2
    github.com/jackc/pgx/v5 v5.10.0
    github.com/go-playground/validator/v10 v10.30.3
    github.com/redis/go-redis/v9 v9.22.0
    github.com/google/uuid v1.6.0
    github.com/joho/godotenv v1.5.1
    github.com/sirupsen/slog v1.9.2
    github.com/jwt-kit/jwt/v5 v5.0.0
)
```

## Build & Run
```
go build ./...          # Compile
go run main.go          # Start server (port 8080)
go vet ./...            # Check for issues
```

## API Base URL
- Development: `http://localhost:8080`
- Android Emulator: `http://10.0.2.2:8080` (localhost from emu perspective)
- Physical Device: `http://10.0.0.1:8080` (or device's local IP)
- Production: `https://your-domain.com`