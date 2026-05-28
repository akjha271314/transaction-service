package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"transaction-service/internal/service"
)

type TransactionHandler struct {
	svc service.TransactionService
}

func NewTransactionHandler(svc service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int64   `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

// Create godoc
// @Summary      Create a transaction
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request body CreateTransactionRequest true "Transaction details"
// @Success      201 {object} models.Transaction
// @Failure      400 {string} string "invalid request body"
// @Failure      422 {string} string "account not found / operation type not found / insufficient credit limit"
// @Security     ApiKeyAuth
// @Router       /transactions [post]
func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := h.svc.CreateTransaction(req.AccountID, req.OperationTypeID, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAccount),
			errors.Is(err, service.ErrInvalidOperationType),
			errors.Is(err, service.ErrInsufficientCredit):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tx)
}
