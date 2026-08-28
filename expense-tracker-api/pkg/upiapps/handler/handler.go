package handler

import (
	"net/http"

	"expense-tracker-api/pkg/response"
	"expense-tracker-api/pkg/upiapps/repository"
)

type UPIAppHandler struct {
	Repo repository.UPIAppRepository
}

func NewUPIAppHandler(repo repository.UPIAppRepository) *UPIAppHandler {
	return &UPIAppHandler{Repo: repo}
}

func (h *UPIAppHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	apps, err := h.Repo.ListAll(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list UPI apps")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    apps,
	})
}
