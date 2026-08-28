package middleware

import (
	"context"
	"net/http"
	"strings"

	"expense-tracker-api/config"
	"expense-tracker-api/pkg/jwt"
	"expense-tracker-api/pkg/response"
	"expense-tracker-api/pkg/users/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const UserKey contextKey = "user"
const RequestIDKey contextKey = "request_id"

type AuthMiddleware struct {
	DB           *pgxpool.Pool
	TokenManager *jwt.TokenManager
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestIDKey, uuid.NewString())

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		// Parse Bearer token: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			response.Error(w, http.StatusUnauthorized, "invalid token format")
			return
		}

		if tokenString == "" {
			response.Error(w, http.StatusUnauthorized, "invalid token format")
			return
		}

		// Validate JWT token
		claims, err := m.TokenManager.ValidateToken(tokenString)
		if err != nil {
			if err == jwt.ErrExpiredToken {
				response.Error(w, http.StatusUnauthorized, "token expired")
			} else {
				response.Error(w, http.StatusUnauthorized, "invalid token")
			}
			return
		}

		// Look up user in database
		repo := repository.NewUserRepository(m.DB)
		user, err := repo.GetByID(ctx, claims.UserID)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "user not found")
			return
		}

		ctx = context.WithValue(ctx, UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewAuthMiddleware(db *pgxpool.Pool) *AuthMiddleware {
	cfg := config.LoadConfig()
	return &AuthMiddleware{
		DB:           db,
		TokenManager: cfg.TokenManager,
	}
}
