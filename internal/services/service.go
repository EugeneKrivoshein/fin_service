package service

import (
	"context"

	repo "github.com/EugeneKrivoshein/fin_service/internal/postgres"
)

// Service содержит бизнес-логику
type Service struct {
	repo *repo.Repository
}

// NewService создаёт новый экземпляр Service
func NewService(r *repo.Repository) *Service {
	return &Service{
		repo: r,
	}
}

// Deposit пополняет баланс пользователя
func (s *Service) Deposit(ctx context.Context, userID int64, amount float64) error {
	return s.repo.Deposit(ctx, userID, amount)
}

// Transfer переводит деньги от одного пользователя к другому
func (s *Service) Transfer(ctx context.Context, senderID, receiverID int64, amount float64) error {
	return s.repo.Transfer(ctx, senderID, receiverID, amount)
}

// GetTransactions возвращает 10 последних транзакций пользователя
func (s *Service) GetTransactions(ctx context.Context, userID int64) ([]repo.Transaction, error) {
	return s.repo.GetTransactions(ctx, userID)
}
