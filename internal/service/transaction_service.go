package service

import (
	"database/sql"
	"errors"
	"math"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

var ErrInvalidAccount = errors.New("account not found")
var ErrInvalidOperationType = errors.New("operation type not found")
var ErrInsufficientBalance = errors.New("insufficient balance")

type TransactionService interface {
	CreateTransaction(accountID, operationTypeID int64, amount float64) (*models.Transaction, error)
}

type transactionService struct {
	txRepo      repository.TransactionRepository
	accountRepo repository.AccountRepository
	txRunner    repository.TxRunner
}

func NewTransactionService(
	txRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	txRunner repository.TxRunner,
) TransactionService {
	return &transactionService{txRepo: txRepo, accountRepo: accountRepo, txRunner: txRunner}
}

func (s *transactionService) CreateTransaction(accountID, operationTypeID int64, amount float64) (*models.Transaction, error) {
	if _, err := s.accountRepo.FindByID(accountID); err != nil {
		return nil, ErrInvalidAccount
	}

	opType, err := s.txRepo.FindOperationType(operationTypeID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidOperationType
		}
		return nil, err
	}

	signedAmount := applySign(opType.IsCredit, amount)

	var result *models.Transaction
	err = s.txRunner.RunInTx(func(tx *sql.Tx) error {
		if err := s.accountRepo.UpdateBalanceTx(tx, accountID, signedAmount); err != nil {
			if errors.Is(err, repository.ErrInsufficientBalance) {
				return ErrInsufficientBalance
			}
			return err
		}
		result, err = s.txRepo.CreateTx(tx, accountID, operationTypeID, signedAmount)
		return err
	})

	return result, err
}

func applySign(isCredit bool, amount float64) float64 {
	if isCredit {
		return math.Abs(amount)
	}
	return -math.Abs(amount)
}
