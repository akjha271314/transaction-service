package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"transaction-service/internal/repository"
	"transaction-service/internal/service"
)

type AccountHandler struct {
	svc service.AccountService
}

func NewAccountHandler(svc service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

type createAccountRequest struct {
	DocumentNumber string  `json:"document_number"`
	CreditLimit    float64 `json:"credit_limit"`
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.DocumentNumber == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	account, err := h.svc.CreateAccount(req.DocumentNumber, req.CreditLimit)
	if err != nil {
		http.Error(w, "could not create account", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("accountId"), 10, 64)
	if err != nil {
		http.Error(w, "invalid account id", http.StatusBadRequest)
		return
	}

	account, err := h.svc.GetAccount(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "account not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
