package handler

import (
	"encoding/json"
	"net/http"

	"expense-tracker-api/internal/middleware"
	"expense-tracker-api/pkg/response"
	"expense-tracker-api/pkg/transactions/model"
	"expense-tracker-api/pkg/transactions/repository"
	usermodel "expense-tracker-api/pkg/users/model"
)

type TransactionHandler struct {
	Repo repository.TransactionRepository
}

func NewTransactionHandler(repo repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{Repo: repo}
}

type CreateTransactionRequest struct {
	FamilyID      string  `json:"family_id"`
	Type          string  `json:"type"`
	Category      string  `json:"category"`
	PaymentMethod string  `json:"payment_method"`
	UPIApp        string  `json:"upi_app,omitempty"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
	Description   string  `json:"description,omitempty"`
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserKey).(*usermodel.User)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FamilyID == "" {
		response.Error(w, http.StatusBadRequest, "family_id is required")
		return
	}
	if req.Type == "" || (req.Type != "DEBIT" && req.Type != "CREDIT") {
		response.Error(w, http.StatusBadRequest, "type must be DEBIT or CREDIT")
		return
	}
	if req.Amount <= 0 {
		response.Error(w, http.StatusBadRequest, "amount must be positive")
		return
	}
	if req.Date == "" {
		response.Error(w, http.StatusBadRequest, "date is required")
		return
	}

	tx := &model.Transaction{
		FamilyID:      req.FamilyID,
		UserID:        &user.ID,
		Type:          req.Type,
		Category:      &req.Category,
		PaymentMethod: &req.PaymentMethod,
		UPIApp:        &req.UPIApp,
		Amount:        req.Amount,
		Date:          req.Date,
		Description:   req.Description,
	}

	if err := h.Repo.Create(ctx, tx); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create transaction")
		return
	}

	response.JSON(w, http.StatusCreated, response.HTTPResponse{
		Success: true,
		Data:    tx,
		Message: "Transaction created successfully",
	})
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, ok := ctx.Value(middleware.UserKey).(*usermodel.User)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	txID := r.URL.Query().Get("id")
	if txID == "" {
		response.Error(w, http.StatusBadRequest, "transaction id is required")
		return
	}

	tx, err := h.Repo.GetByID(ctx, txID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "transaction not found")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    tx,
	})
}

func (h *TransactionHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
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

	txs, err := h.Repo.ListByUser(ctx, user.ID, &familyID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list transactions")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    txs,
	})
}

func (h *TransactionHandler) ListMyTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(middleware.UserKey).(*usermodel.User)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	txs, err := h.Repo.ListByUser(ctx, user.ID, nil)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list transactions")
		return
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    txs,
	})
}
