package models

type Account struct {
	ID             int64   `json:"account_id"`
	DocumentNumber string  `json:"document_number"`
	Balance        float64 `json:"balance"`
}
