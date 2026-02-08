package services

import (
	"kasir-api/model"
	"kasir-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
    return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []model.CheckoutItem, useLock bool) (*model.Transaction, error) {
		return s.repo.CreateTransaction(items)
}

func (s *TransactionService) GetReport(startDate, endDate string) (*model.ReportResponse, error) {
    // Service meneruskan permintaan ke Repository
    return s.repo.GetReport(startDate, endDate)
}