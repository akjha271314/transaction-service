package service

import (
	"errors"
	"math"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

var ErrInvalidAccount = errors.New("account not found")
var ErrInvalidOperationType = errors.New("operation type not found")

// creditVoucherID is the only operation type stored with a positive amount.
const creditVoucherID = 4

type TransactionService interface {
	CreateTransaction(accountID, operationTypeID int64, amount float64) (*models.Transaction, error)
}

type transactionService struct {
	txRepo      repository.TransactionRepository
	accountRepo repository.AccountRepository
}

func NewTransactionService(txRepo repository.TransactionRepository, accountRepo repository.AccountRepository) TransactionService {
	return &transactionService{txRepo: txRepo, accountRepo: accountRepo}
}

func (s *transactionService) CreateTransaction(accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
	if _, err := s.accountRepo.FindByID(accountID); err != nil {
		return nil, ErrInvalidAccount
	}

	exists, err := s.txRepo.OperationTypeExists(operationTypeID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrInvalidOperationType
	}

	return s.txRepo.Create(accountID, operationTypeID, applySign(operationTypeID, amount))
}

func applySign(operationTypeID int64, amount float64) float64 {
	if operationTypeID == creditVoucherID {
		return math.Abs(amount)
	}
	return -math.Abs(amount)
}
