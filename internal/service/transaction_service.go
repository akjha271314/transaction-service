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

// creditVoucherID is the only operation type stored with a positive amount.
const creditVoucherID = 4

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

	exists, err := s.txRepo.OperationTypeExists(operationTypeID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrInvalidOperationType
	}

	signedAmount := applySign(operationTypeID, amount)

	var result *models.Transaction
	err = s.txRunner.RunInTx(func(tx *sql.Tx) error {
		if err := s.accountRepo.UpdateBalanceTx(tx, accountID, signedAmount); err != nil {
			return err
		}
		result, err = s.txRepo.CreateTx(tx, accountID, operationTypeID, signedAmount)
		return err
	})

	return result, err
}

func applySign(operationTypeID int64, amount float64) float64 {
	if operationTypeID == creditVoucherID {
		return math.Abs(amount)
	}
	return -math.Abs(amount)
}
