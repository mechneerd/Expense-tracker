# Expense Tracker

Full-stack expense tracking application with family support.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.23 + Chi Router |
| Database | PostgreSQL 15 |
| Cache | Redis 7 |
| ORM | sqlc (type-safe queries) |
| Mobile | Kotlin + Android |

---

## Project Structure

```
Expense/
├── expense-tracker-api/       # Go Backend
│   ├── main.go               # Entry point
│   ├── config/               # Configuration
│   ├── internal/middleware/   # Auth middleware
│   ├── pkg/                  # Domain packages
│   │   ├── handler/          # User handler
│   │   ├── jwt/              # JWT tokens
│   │   ├── users/            # User domain
│   │   ├── families/         # Family domain
│   │   ├── transactions/     # Transaction domain
│   │   ├── categories/       # Categories
│   │   ├── paymentmethods/   # Payment methods
│   │   └── upiapps/          # UPI apps
│   ├── queries/              # SQL queries
│   ├── migrations/           # Database schema
│   ├── Dockerfile            # Docker build
│   └── docker-compose.yml    # Local dev
│
├── mobile-app/               # Android App
│   └── app/                  # Kotlin source
│
└── README.md                 # This file
```

---

## API Endpoints

### Auth
```
POST /api/v1/auth/google      # Register/Login
POST /api/v1/auth/verify-otp   # Verify OTP
POST /api/v1/auth/refresh      # Refresh token
POST /api/v1/auth/logout       # Logout
```

### Users
```
GET    /api/v1/users/me   # Get profile
PATCH  /api/v1/users/me   # Update profile
```

### Families
```
POST   /api/v1/families/      # Create family
GET    /api/v1/families/      # Get family
GET    /api/v1/families/me    # List members
POST   /api/v1/families/invite # Invite member
```

### Transactions
```
POST /api/v1/transactions/     # Create transaction
GET  /api/v1/transactions/     # List family transactions
GET  /api/v1/transactions/me   # List my transactions
```

### Reference Data
```
GET /api/v1/categories/        # List categories
GET /api/v1/payment-methods/   # List payment methods
GET /api/v1/upi-apps/          # List UPI apps
```

### Health
```
GET /health                    # Health check
```

---

## Local Development

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- PostgreSQL 15+ (or use Docker)
- Redis 7+ (or use Docker)

### Quick Start

```bash
# 1. Clone
git clone <repo-url>
cd Expense/expense-tracker-api

# 2. Copy env file
cp .env.example .env

# 3. Start database and Redis
docker-compose up -d postgres redis

# 4. Run the API
go run main.go
```

API runs on `http://localhost:8080`

### Run Tests
```bash
go test ./...
```

### Build Binary
```bash
go build -o expense-api .
```

---

## Free Deployment (Railway + Supabase)

### Step 1: Database (Supabase - Free)

1. Go to [supabase.com](https://supabase.com)
2. Sign up with GitHub
3. Create new project:
   - Name: `expense-tracker`
   - Database password: (save this)
   - Region: closest to you
4. Go to **SQL Editor**
5. Paste contents of `migrations/001_initial_schema.up.sql`
6. Click **Run**
7. Go to **Settings > Database**
8. Copy the **Connection string > URI**
   - Format: `postgresql://postgres.[ref]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres`

### Step 2: Backend (Railway - Free)

1. Go to [railway.app](https://railway.app)
2. Sign up with GitHub
3. Click **New Project > Deploy from GitHub repo**
4. Select your `Expense` repository
5. Select `expense-tracker-api` as root directory
6. Go to **Variables** tab, add:

```
DATABASE_URL=postgresql://postgres.[ref]:[password]@aws-0-[region].pooler.supabase.com:6543/postgres
JWT_SECRET_KEY=your-random-secret-key-here-make-it-long
REDIS_URL=redis://default:[password]@redis.railway.internal:6379
```

7. Go to **Settings**:
   - Build Command: `go build -o main .`
   - Start Command: `./main`
   - Port: `8080`

8. Railway auto-deploys. Get your URL from **Settings > Networking**

### Step 3: Redis (Upstash - Free)

1. Go to [upstash.com](https://upstash.com)
2. Sign up with GitHub
3. Click **Create Database**
   - Name: `expense-redis`
   - Type: **Regional**
   - Region: same as Railway
4. Copy the **Redis URL**
5. Add to Railway variables:
   ```
   REDIS_URL=redis://default:[password]@.[region].upstash.io:[port]
   ```

### Alternative: All-in-One on Railway

Railway has built-in PostgreSQL and Redis:

1. Create new Railway project
2. Add **PostgreSQL** service (free $5 credit/month)
3. Add **Redis** service (free $5 credit/month)
4. Add your Go app service
5. Railway auto-links the variables:
   - `PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE`
   - `REDIS_URL`
6. Update your code to use these variables or set `DATABASE_URL` manually

---

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_URL` | Yes | `postgres://postgres:postgres@localhost:5432/expense_tracker?sslmode=disable` | PostgreSQL connection |
| `JWT_SECRET_KEY` | Yes | `your-default-secret-key` | JWT signing secret (CHANGE THIS!) |
| `REDIS_URL` | No | `localhost:6379` | Redis connection |
| `PORT` | No | `8080` | Server port |
| `OT_LIFETIME_MINUTES` | No | `5` | OTP expiry |

---

## Docker

### Build and Run
```bash
# Build image
docker build -t expense-api .

# Run with Docker Compose (includes PostgreSQL + Redis)
docker-compose up --build

# Run only the API (external DB)
docker run -p 8080:8080 \
  -e DATABASE_URL=your-db-url \
  -e JWT_SECRET_KEY=your-secret \
  expense-api
```

---

## Mobile App

### Setup
1. Open `mobile-app/` in Android Studio
2. Update API base URL in your network config
3. Build and run on emulator or device

### API Configuration
- Emulator: `http://10.0.2.2:8080`
- Device: `http://<your-computer-ip>:8080`
- Production: `https://your-railway-app.up.railway.app`

---

## Database Schema

Tables created by `migrations/001_initial_schema.up.sql`:

| Table | Purpose |
|-------|---------|
| `users` | User accounts |
| `families` | Family groups |
| `family_members` | Family membership |
| `family_invitations` | Invitations |
| `transactions` | Income/Expense records |
| `transaction_categories` | Categories |
| `payment_methods` | Payment methods |
| `upi_apps` | UPI applications |
| `otp_verifications` | OTP tracking |
| `audit_logs` | Audit trail |

---

## Development Commands

```bash
# Run API
go run main.go

# Run tests
go test ./...

# Build binary
go build -o expense-api .

# Generate sqlc code
sqlc generate

# Start dev services
docker-compose up -d

# View logs
docker-compose logs -f api
```

---

## License

MIT
