# Expense Tracker - Suggestions & Improvements

Based on a deep analysis of the codebase, here are prioritized suggestions for security, features, and architecture improvements.

---

## 🔒 Security Issues

### 1. JWT Token Validation Not Implemented (`internal/middleware/auth.go:28`)
```go
_ = authHeader // TODO: Validate JWT token
```
**Risk**: Authentication is effectively disabled - any request passes through without JWT validation.
**Fix**: Implement proper JWT validation using `github.com/jwt-kit/jwt/v5`:
- Parse the authorization header
- Validate token signature against `JWT_SECRET_KEY` from config
- Extract user claims and set in context
- Return 401 for invalid/missing tokens

### 2. Passwords Not Hashed
- User model stores `google_id` but no password field exists
- If password authentication is added later, must use bcrypt or similar
**Fix**: Add `password_hash` field to user model; never store plaintext passwords

### 3. No Rate Limiting on OTP Verification
- OTP endpoint has no rate limiting (3 requests/10 min mentioned in docs but not enforced)
**Fix**: Add Redis-based rate limiting similar to the documented design

### 4. CORS Not Configured
- Chi router has no CORS middleware
**Fix**: Add `github.com/go-chi/cors` with appropriate origins for mobile app

### 5. No Input Validation Beyond struct tags
- Handlers read form values directly without validation
**Fix**: Use `github.com/go-playder/validator/v10` struct validation or ensure all inputs are sanitized

---

## ⚡ Features to Implement

### High Priority (Business Logic)

| # | Feature | Status | Location |
|---|---------|--------|----------|
| 1 | **Full JWT token generation** | ❌ Stub | `handler/user.go:101-103` |
| 2 | **OTP verification with Redis** | ❌ Stub | `handler/user.go:61-67` |
| 3 | **Refresh token rotation** | ❌ Stub | `handler/user.go:69-75` |
| 4 | **Transaction creation with DB** | ❌ Stub | `pkg/transactions/repository.go:22-25` |
| 5 | **Category repository SQL queries** | ❌ Stub | `pkg/categories/repository.go:20-28` |
| 6 | **Payment method repository SQL queries** | ❌ Stub | `pkg/paymentmethods/repository.go:20-28` |
| 7 | **UPI app repository SQL queries** | ❌ Stub | `pkg/upiapps/repository.go:20-28` |
| 8 | **Family member management** | ❌ Stubs | `handler/family.go` |
| 9 | **Invitation system** | ❌ Stubs | `handler/invitations.go` |
| 10 | **Role-based access control** | ⚠️ Partial | `middleware/auth.go` |

### Medium Priority (User Experience)

| # | Feature | Notes |
|---|---------|-------|
| 11 | **Dashboard with period filter** | Already routed; needs `dashboardHandler.GetDashboard` implementation |
| 12 | **Transaction filtering by category/date** | Requires repository query enhancements |
| 13 | **Family joining by unique code** | Missing endpoint + UI flow |
| 14 | **Profile image upload** | AvatarURL exists but no upload endpoint |
| 15 | **Email verification flow** | OTP system exists but email sending is stubbed |

### Low Priority (Polish)

| # | Feature |
|---|---------|
| 16 | **API documentation** (OpenAPI/Swagger) |
| 17 | **Request/response logging** (correlation IDs already partially done) |
| 18 | **Graceful shutdown** handler |
| 19 | **Test suite** (currently no test files) |
| 20 | **Configuration validation** at startup |

---

## 🏗️ Architecture Improvements

### 1. Generate sqlc Queries
```bash
sqlc generate
```
**Benefit**: Type-safe SQL queries instead of raw `db.QueryContext` calls; eliminates N+1 query problems and SQL injection risks.

### 2. Repository Pattern Consistency
Currently mixed implementations:
- `pkg/users/repository` - full SQL implementation
- `pkg/categories/repository` - interface only with TODO comments
- `pkg/paymentmethods/repository` - interface only with TODO comments

**Fix**: Either implement all repositories with sqlc or ensure consistent pattern across all modules.

### 3. Move Shared Models to `pkg/`
Several models duplicated or in inconsistent locations:
- `pkg/users/model/user.go` ✅
- `pkg/categories/model/categories.go` ✅ (includes PaymentMethod, UPIApp)
- `internal/` has no models currently

**Suggestion**: Keep all models in `pkg/` and have `internal/` modules depend on `pkg/` interfaces.

### 4. Add Middleware for Request ID Correlation
Currently `internal/middleware/auth.go` sets a hardcoded `"req-0000"`:
```go
ctx = context.WithValue(ctx, RequestIDKey, "req-0000")
```
**Fix**: Generate unique request IDs per request using `github.com/google/uuid` or `github.com/sony/gtrace`.

### 5. Implement Proper Error Handling
Many handlers return 200 OK even on errors:
```go
response.JSON(w, http.StatusOK, response.HTTPResponse{
    Success: true,
    Message: "Transaction created",
})
```
**Fix**: Return appropriate HTTP status codes (400, 404, 422, 500) based on error type.

### 6. Add Configuration Validation
`config.LoadConfig()` uses defaults if env vars missing:
```go
JWTSecretKey: getEnv("JWT_SECRET_KEY", "your-default-secret-key"),
```
**Fix**: Validate critical configs at startup and exit if missing:
- DATABASE_URL must be set
- JWT_SECRET_KEY must not be the default value
- REDIS_ADDR must be reachable

---

## 📊 Database & Performance

### 1. Add Soft Delete Queries
Transaction model has `deleted_at` but no queries filter by it:
- `transactionRepository.ListByFamily` and `ListByUser` should exclude soft-deleted records

### 2. Add Database Connection Pool Settings
Current `sql.Open` uses default pgx pool settings. Consider:
```go
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### 3. Add Auditing for All Writes
`audit_logs` table exists but no code populates it. Consider automatic audit logging on:
- User creation/update
- Family member changes
- Transaction creation

---

## 📱 Mobile App Considerations

### 1. API Client Security
- Add request timeout to Retrofit client
- Implement certificate pinning for production
- Securely store JWT in `SharedPreferences` (currently stored as plain text)

### 2. Network Configuration
- Emulator: `http://10.0.2.2:8080` ✅
- Physical device: needs IP detection logic
- Add fallback to localhost for testing

### 3. Offline Support
- Currently no offline/cache capability
- Consider Room database for transaction caching

---

## 🧪 Testing Recommendations

### 1. Add Integration Tests
```go
go test ./...  # Currently finds nothing
```

### 2. Test Critical Flows
- Google login → OTP verification → API calls → Transaction creation
- Family creation → Member invitation → Acceptance → Transaction under family

### 3. Use Testcontainers or SQLC Mocks
- For PostgreSQL integration tests
- Or generate sqlc mocks: `sqlc generate --mock`

---

## 📝 Implementation Priority Order

| Priority | Item | Effort | Impact |
|----------|------|--------|--------|
| P0 | Implement JWT validation in middleware | 1 day | 🔴 Critical |
| P0 | Implement OTP verification + Redis | 2 days | 🔴 Critical |
| P0 | Implement transaction creation with DB | 2 days | 🟠 High |
| P1 | Generate sqlc queries for all modules | 1 day | 🟠 High |
| P1 | Implement refresh token logic | 1 day | 🟠 High |
| P2 | Add rate limiting to auth endpoints | 0.5 day | 🟡 Medium |
| P2 | Add CORS middleware | 0.5 day | 🟡 Medium |
| P3 | Add comprehensive test suite | 3 days | 🟢 Low |
| P3 | Add API documentation (OpenAPI) | 2 days | 🟢 Low |

---