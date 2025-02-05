package service

import (
	"context"

	repo "github.com/EugeneKrivoshein/fin_service/internal/postgres"
)

type Service struct {
	repo repo.Repository
}

func NewService(r repo.Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) Deposit(ctx context.Context, userID int64, amount float64) error {
	return s.repo.Deposit(ctx, userID, amount)
}

func (s *Service) Transfer(ctx context.Context, senderID, receiverID int64, amount float64) error {
	return s.repo.Transfer(ctx, senderID, receiverID, amount)
}

func (s *Service) GetTransactions(ctx context.Context, userID int64) ([]repo.Transaction, error) {
	return s.repo.GetTransactions(ctx, userID)
}
