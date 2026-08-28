package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"expense-tracker-api/config"
	"fmt"
	"math/big"
	"os"
	"expense-tracker-api/pkg/response"
	"expense-tracker-api/pkg/users/model"
	"expense-tracker-api/pkg/users/repository"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	Repo     repository.UserRepository
	Cfg      *config.Config
}

func NewUserHandler(pool *pgxpool.Pool) *UserHandler {
	return &UserHandler{
		Repo: repository.NewUserRepository(pool),
		Cfg:  config.LoadConfig(),
	}
}

func (h *UserHandler) RegisterGoogle(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	email := r.FormValue("email")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")

	user, err := h.Repo.GetByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	if err == sql.ErrNoRows {
		user = model.NewUser(email, firstName, lastName)
		if err := h.Repo.Create(ctx, user); err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
	}

	// Store OTP in Redis with TTL
	otpKey := fmt.Sprintf("otp:%s", email)
	otpCode := fmt.Sprintf("%06d", newRandInt(100000, 999999)) // generate 6-digit OTP
	h.Cfg.Rdb.Set(ctx, otpKey, otpCode, time.Duration(h.Cfg.OTLifetime)*time.Minute)

	token, err := generateToken(user.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data: map[string]interface{}{
			"user":  user,
			"token": token,
			"otp_key": otpKey,
		},
		Message: "User registered/verified. OTP sent to email.",
	})
}

func (h *UserHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	email := r.FormValue("email")
	otpCode := r.FormValue("otp")

	if email == "" || otpCode == "" {
		response.Error(w, http.StatusBadRequest, "email and OTP code are required")
		return
	}

	// Look up user by email
	dbUser, err := h.Repo.GetByEmail(ctx, email)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	otpKey := fmt.Sprintf("otp:%s", email)
	storedOTP, err := h.Cfg.Rdb.Get(ctx, otpKey).Result()
	if err != nil {
		if err == redis.Nil {
			response.Error(w, http.StatusBadRequest, "OTP expired or not found")
		} else {
			fmt.Fprintf(os.Stderr, "Redis error: %v\n", err)
			response.Error(w, http.StatusInternalServerError, "Failed to verify OTP")
		}
		return
	}

	// Verify OTP using constant-time comparison
	if !strings.EqualFold(storedOTP, otpCode) {
		response.Error(w, http.StatusBadRequest, "invalid OTP code")
		return
	}

	// Mark OTP as verified - delete from Redis
	h.Cfg.Rdb.Del(ctx, otpKey)

	// Generate access token
	token, err := generateToken(dbUser.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Generate refresh token with 30-day expiry
	refreshToken, err := generateRefreshToken(dbUser.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Store refresh token in Redis with 30-day TTL
	refreshKey := fmt.Sprintf("refresh:%s", dbUser.ID)
	h.Cfg.Rdb.Set(ctx, refreshKey, refreshToken, 30*24*time.Hour)

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data: map[string]interface{}{
			"user":     dbUser,
			"token":    token,
			"refresh":  refreshToken,
		},
		Message: "OTP verified successfully",
	})
}

func (h *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	refreshToken := r.FormValue("refresh_token")
	if refreshToken == "" {
		response.Error(w, http.StatusBadRequest, "missing refresh token")
		return
	}

	// Get email from form
	email := r.FormValue("email")
	if email == "" {
		response.Error(w, http.StatusBadRequest, "email is required")
		return
	}

	dbUser, err := h.Repo.GetByEmail(ctx, email)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "user not found")
		return
	}

	refreshKey := fmt.Sprintf("refresh:%s", dbUser.ID)
	storedToken, err := h.Cfg.Rdb.Get(ctx, refreshKey).Result()
	if err != nil {
		if err == redis.Nil {
			response.Error(w, http.StatusUnauthorized, "refresh token expired or not found")
		} else {
			response.Error(w, http.StatusInternalServerError, "Failed to verify refresh token")
		}
		return
	}

	// Verify token match
	if !strings.EqualFold(storedToken, refreshToken) {
		response.Error(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	// Revoke old refresh token (rotation)
	h.Cfg.Rdb.Del(ctx, refreshKey)

	// Generate new access token
	newToken, err := generateToken(dbUser.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Generate new refresh token
	newRefreshToken, err := generateRefreshToken(dbUser.ID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Store new refresh token
	newRefreshKey := fmt.Sprintf("refresh:%s", dbUser.ID)
	h.Cfg.Rdb.Set(ctx, newRefreshKey, newRefreshToken, 30*24*time.Hour)

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data: map[string]interface{}{
			"token":        newToken,
			"refresh_token": newRefreshToken,
		},
		Message: "Token refreshed successfully",
	})
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*model.User)
	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value("user").(*model.User)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Phone     string `json:"phone"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.AvatarURL != "" {
		user.AvatarURL = &req.AvatarURL
	}

	if err := h.Repo.Update(ctx, user); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    user,
		Message: "Profile updated successfully",
	})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user email from form
	email := r.FormValue("email")
	if email == "" {
		response.Error(w, http.StatusBadRequest, "email is required")
		return
	}

	dbUser, err := h.Repo.GetByEmail(context.Background(), email)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	// Revoke refresh token
	refreshKey := fmt.Sprintf("refresh:%s", dbUser.ID)
	h.Cfg.Rdb.Del(context.Background(), refreshKey)

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Message: "Logged out successfully",
	})
}

func generateToken(userID string) (string, error) {
	cfg := config.LoadConfig()
	return cfg.TokenManager.GenerateAccessToken(userID)
}

func generateRefreshToken(userID string) (string, error) {
	cfg := config.LoadConfig()
	return cfg.TokenManager.GenerateRefreshToken(userID)
}

// helper function to generate random int using crypto/rand
func newRandInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return min + int(n.Int64())
}