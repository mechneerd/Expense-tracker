package handler

import (
	"net/http"

	"expense-tracker-api/pkg/categories/repository"
	"expense-tracker-api/pkg/response"
)

type CategoryHandler struct {
	Repo repository.CategoryRepository
}

func NewCategoryHandler(repo repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{Repo: repo}
}

func (h *CategoryHandler) ListByType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categoryType := r.URL.Query().Get("type")
	if categoryType == "" {
		categoryType = "DEBIT"
	}

	categories, err := h.Repo.ListByType(ctx, categoryType)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list categories")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    categories,
	})
}
