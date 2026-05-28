package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"transaction-service/internal/service"
)

type AccountHandler struct {
	svc service.AccountService
}

func NewAccountHandler(svc service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

type CreateAccountRequest struct {
	DocumentNumber string  `json:"document_number"`
	Balance        float64 `json:"balance"`
}

// Create godoc
// @Summary      Create an account
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        request body CreateAccountRequest true "Account details"
// @Success      201 {object} models.Account
// @Failure      400 {string} string "invalid request body / balance must not be negative"
// @Failure      409 {string} string "account with this document number already exists"
// @Security     ApiKeyAuth
// @Router       /accounts [post]
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DocumentNumber == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Balance < 0 {
		http.Error(w, "balance must not be negative", http.StatusBadRequest)
		return
	}

	account, err := h.svc.CreateAccount(req.DocumentNumber, req.Balance)
	if err != nil {
		if errors.Is(err, service.ErrDuplicateAccount) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, "could not create account", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// GetByID godoc
// @Summary      Get an account
// @Tags         accounts
// @Produce      json
// @Param        accountId path int true "Account ID"
// @Success      200 {object} models.Account
// @Failure      400 {string} string "invalid account id"
// @Failure      404 {string} string "account not found"
// @Security     ApiKeyAuth
// @Router       /accounts/{accountId} [get]
func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("accountId"), 10, 64)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	account, err := h.svc.GetAccount(id)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
