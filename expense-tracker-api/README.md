# Expense Tracker API

A full-stack expense tracking application with family support, built with Go and PostgreSQL.

## Features

- **Authentication**: OTP-based auth with JWT tokens and refresh token rotation
- **Family Management**: Create families, invite members, track expenses by family
- **Transactions**: Record income/expense with categories, payment methods, and UPI apps
- **Dashboard**: Financial summaries with period-based filtering
- **Type-Safe Queries**: Uses sqlc for compile-time verified SQL queries

## Tech Stack

- **Backend**: Go 1.23+, Chi Router, pgx (PostgreSQL)
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Code Generation**: sqlc for type-safe database queries
- **Testing**: testify

## Quick Start

### Prerequisites

- Go 1.23+
- PostgreSQL 15+
- Redis 7+

### Local Development

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd expense-tracker-api
   ```

2. Copy environment file:
   ```bash
   cp .env.example .env
   ```

3. Start services with Docker:
   ```bash
   docker-compose up -d postgres redis
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

### Docker

Build and run with Docker Compose:
```bash
docker-compose up --build
```

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/google` | Register/login with Google |
| POST | `/api/v1/auth/verify-otp` | Verify OTP code |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/logout` | Logout (revoke refresh token) |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users/me` | Get current user |
| PATCH | `/api/v1/users/me` | Update profile |

### Families
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/families/` | Create family |
| GET | `/api/v1/families/` | Get family by ID |
| GET | `/api/v1/families/me` | List family members |
| POST | `/api/v1/families/invite` | Invite member |

### Transactions
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/transactions/` | Create transaction |
| GET | `/api/v1/transactions/` | List family transactions |
| GET | `/api/v1/transactions/me` | List my transactions |

### Reference Data
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/categories/` | List categories by type |
| GET | `/api/v1/payment-methods/` | List payment methods |
| GET | `/api/v1/upi-apps/` | List UPI apps |

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |

## Project Structure

```
expense-tracker-api/
├── main.go                    # Entry point
├── config/                    # Configuration
├── internal/
│   ├── middleware/            # Auth middleware
│   └── db/
│       ├── generated/        # sqlc generated code
│       └── dbutil/           # Database utilities
├── pkg/
│   ├── handler/              # User handler
│   ├── jwt/                  # JWT token manager
│   ├── response/             # HTTP response helpers
│   ├── users/                # User domain
│   ├── families/             # Family domain
│   ├── transactions/         # Transaction domain
│   ├── categories/           # Category domain
│   ├── paymentmethods/       # Payment method domain
│   └── upiapps/              # UPI app domain
├── migrations/               # Database migrations
└── queries/                  # SQL query definitions
```

## Development

### Run Tests
```bash
go test ./...
```

### Build
```bash
go build -o expense-tracker-api .
```

### Generate sqlc Code
```bash
sqlc generate
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/expense_tracker?sslmode=disable` | PostgreSQL connection URL |
| `JWT_SECRET_KEY` | `your-default-secret-key` | JWT signing secret |
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | `` | Redis password |
| `OT_LIFETIME_MINUTES` | `5` | OTP expiry in minutes |
| `PORT` | `8080` | Server port |

## License

MIT
