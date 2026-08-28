package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"expense-tracker-api/internal/middleware"
	"expense-tracker-api/pkg/families/model"
	"expense-tracker-api/pkg/families/repository"
	"expense-tracker-api/pkg/response"
	usermodel "expense-tracker-api/pkg/users/model"
)

type FamilyHandler struct {
	Repo repository.FamilyRepository
}

func NewFamilyHandler(repo repository.FamilyRepository) *FamilyHandler {
	return &FamilyHandler{Repo: repo}
}

type CreateFamilyRequest struct {
	Name string `json:"name"`
}

func (h *FamilyHandler) CreateFamily(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserKey).(*usermodel.User)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateFamilyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "family name is required")
		return
	}

	uniqueCode := generateUniqueCode()

	family := &model.Family{
		Name:       req.Name,
		UniqueCode: uniqueCode,
		CreatedBy:  user.ID,
		Status:     "active",
	}

	if err := h.Repo.Create(ctx, family); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create family")
		return
	}

	member := &model.FamilyMember{
		FamilyID:   family.ID,
		UserID:     user.ID,
		FamilyRole: "HEAD",
		Status:     "ACTIVE",
		JoinedAt:   time.Now().Format(time.RFC3339),
	}

	if err := h.Repo.AddMember(ctx, member); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to add family member")
		return
	}

	response.JSON(w, http.StatusCreated, response.HTTPResponse{
		Success: true,
		Data:    family,
		Message: "Family created successfully",
	})
}

func (h *FamilyHandler) GetFamily(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserKey).(*usermodel.User)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	familyID := r.URL.Query().Get("family_id")
	if familyID == "" {
		response.Error(w, http.StatusBadRequest, "family_id is required")
		return
	}

	family, err := h.Repo.GetByID(ctx, familyID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "family not found")
		return
	}

	_ = user // User context available for authorization checks

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    family,
	})
}

func (h *FamilyHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := ctx.Value(middleware.UserKey).(*usermodel.User)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	familyID := r.URL.Query().Get("family_id")
	if familyID == "" {
		response.Error(w, http.StatusBadRequest, "family_id is required")
		return
	}

	members, err := h.Repo.ListMembers(ctx, familyID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list members")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    members,
	})
}

func (h *FamilyHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserKey).(*usermodel.User)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req struct {
		FamilyID string `json:"family_id"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FamilyID == "" || req.Email == "" {
		response.Error(w, http.StatusBadRequest, "family_id and email are required")
		return
	}

	members, err := h.Repo.ListMembers(ctx, req.FamilyID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to verify family membership")
		return
	}

	isHead := false
	for _, m := range members {
		if m.UserID == user.ID && m.FamilyRole == "HEAD" {
			isHead = true
			break
		}
	}

	if !isHead {
		response.Error(w, http.StatusForbidden, "only family head can invite members")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Message: "Invitation sent to " + req.Email,
	})
}

func generateUniqueCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}
	return string(code)
}
