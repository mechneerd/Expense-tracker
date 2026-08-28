package handler

import (
	"net/http"

	"expense-tracker-api/pkg/paymentmethods/repository"
	"expense-tracker-api/pkg/response"
)

type PaymentMethodHandler struct {
	Repo repository.PaymentMethodRepository
}

func NewPaymentMethodHandler(repo repository.PaymentMethodRepository) *PaymentMethodHandler {
	return &PaymentMethodHandler{Repo: repo}
}

func (h *PaymentMethodHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	methods, err := h.Repo.ListAll(ctx)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list payment methods")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    methods,
	})
}
