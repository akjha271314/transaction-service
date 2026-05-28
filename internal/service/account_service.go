package service

import (
	"transaction-service/internal/models"
	"transaction-service/internal/repository"
)

type AccountService interface {
	CreateAccount(documentNumber string, creditLimit float64) (*models.Account, error)
	GetAccount(id int64) (*models.Account, error)
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(documentNumber string, creditLimit float64) (*models.Account, error) {
	return s.repo.Create(documentNumber, creditLimit)
}

func (s *accountService) GetAccount(id int64) (*models.Account, error) {
	return s.repo.FindByID(id)
}
