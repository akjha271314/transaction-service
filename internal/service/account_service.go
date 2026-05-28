package service

import (
	"errors"

	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrDuplicateAccount = errors.New("account with this document number already exists")

type AccountService interface {
	CreateAccount(documentNumber string, balance float64) (*models.Account, error)
	GetAccount(id int64) (*models.Account, error)
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(documentNumber string, balance float64) (*models.Account, error) {
	acc, err := s.repo.Create(documentNumber, balance)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateAccount) {
			return nil, ErrDuplicateAccount
		}
		return nil, err
	}
	return acc, nil
}

func (s *accountService) GetAccount(id int64) (*models.Account, error) {
	acc, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return acc, nil
}
