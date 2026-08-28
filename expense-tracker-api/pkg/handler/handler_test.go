package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"expense-tracker-api/pkg/response"
	"expense-tracker-api/pkg/users/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockUserRepo struct {
	users map[string]*model.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*model.User)}
}

func (m *mockUserRepo) GetByEmail(email string) (*model.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) GetByID(id string) (*model.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, nil
}

func (m *mockUserRepo) Create(u *model.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) Update(u *model.User) error {
	m.users[u.ID] = u
	return nil
}

func (m *mockUserRepo) VerifyEmail(id string) error {
	if u, ok := m.users[id]; ok {
		u.EmailVerifiedAt = u.CreatedAt
	}
	return nil
}

func TestGetMe_Success(t *testing.T) {
	user := &model.User{
		ID:        "test-user-123",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/me", nil)
	req = req.WithContext(context.WithValue(req.Context(), "user", user))

	handler := &UserHandler{}
	handler.GetMe(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.HTTPResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestUpdateMe_Unauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("PATCH", "/me", nil)

	handler := &UserHandler{}
	handler.UpdateMe(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var resp response.HTTPResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.False(t, resp.Success)
}
