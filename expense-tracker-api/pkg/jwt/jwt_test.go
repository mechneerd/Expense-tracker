package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndValidateAccessToken(t *testing.T) {
	tm := NewTokenManager("test-secret-key", 15*time.Minute, 30*24*time.Hour)

	userID := "test-user-123"
	token, err := tm.GenerateAccessToken(userID)

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := tm.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, userID, claims.Subject)
}

func TestGenerateAndValidateRefreshToken(t *testing.T) {
	tm := NewTokenManager("test-secret-key", 15*time.Minute, 30*24*time.Hour)

	userID := "test-user-456"
	token, err := tm.GenerateRefreshToken(userID)

	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := tm.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	tm := NewTokenManager("test-secret-key", 15*time.Minute, 30*24*time.Hour)

	_, err := tm.ValidateToken("invalid-token-string")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	tm1 := NewTokenManager("secret-1", 15*time.Minute, 30*24*time.Hour)
	tm2 := NewTokenManager("secret-2", 15*time.Minute, 30*24*time.Hour)

	token, err := tm1.GenerateAccessToken("user-123")
	require.NoError(t, err)

	_, err = tm2.ValidateToken(token)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	tm := NewTokenManager("test-secret-key", -1*time.Minute, 30*24*time.Hour)

	token, err := tm.GenerateAccessToken("user-123")
	require.NoError(t, err)

	_, err = tm.ValidateToken(token)
	assert.Error(t, err)
	assert.Equal(t, ErrExpiredToken, err)
}

func TestTokenManager_DifferentAccessAndRefresh(t *testing.T) {
	tm := NewTokenManager("test-secret-key", 15*time.Minute, 30*24*time.Hour)

	accessClaims := &Claims{
		UserID: "user-1",
		RegisteredClaims: tm.accessClaims(),
	}

	refreshClaims := &Claims{
		UserID: "user-1",
		RegisteredClaims: tm.refreshClaims(),
	}

	assert.WithinDuration(t, time.Now().Add(15*time.Minute), accessClaims.ExpiresAt.Time, time.Minute)
	assert.WithinDuration(t, time.Now().Add(30*24*time.Hour), refreshClaims.ExpiresAt.Time, time.Minute)
}
