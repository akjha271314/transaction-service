package models

type OperationType struct {
	ID          int64  `json:"operation_type_id"`
	Description string `json:"description"`
	IsCredit    bool   `json:"is_credit"`
}
