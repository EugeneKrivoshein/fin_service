package service

import (
	"context"
	"errors"
	"testing"

	"github.com/EugeneKrivoshein/fin_service/internal/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Deposit(ctx context.Context, userID int64, amount float64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

func (m *MockRepository) Transfer(ctx context.Context, senderID, receiverID int64, amount float64) error {
	args := m.Called(ctx, senderID, receiverID, amount)
	return args.Error(0)
}

func (m *MockRepository) GetTransactions(ctx context.Context, userID int64) ([]postgres.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]postgres.Transaction), args.Error(1)
}

func TestDeposit(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	userID := int64(1)
	amount := 100.0

	// Ожидаем, что метод Deposit будет вызван с заданными параметрами и не вернет ошибок
	mockRepo.On("Deposit", mock.Anything, userID, amount).Return(nil)

	err := service.Deposit(context.Background(), userID, amount)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeposit_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Входные данные
	userID := int64(1)
	amount := 100.0

	// Ожидаем, что метод Deposit вызовет ошибку
	mockRepo.On("Deposit", mock.Anything, userID, amount).Return(errors.New("db error"))

	err := service.Deposit(context.Background(), userID, amount)

	// Проверяем, что ошибка возвращена и она соответствует ожидаемой
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestTransfer(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Входные данные
	senderID := int64(1)
	receiverID := int64(2)
	amount := 50.0

	// Ожидаем, что метод Transfer будет вызван с заданными параметрами и не вернет ошибок
	mockRepo.On("Transfer", mock.Anything, senderID, receiverID, amount).Return(nil)

	err := service.Transfer(context.Background(), senderID, receiverID, amount)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestTransfer_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Входные данные
	senderID := int64(1)
	receiverID := int64(2)
	amount := 50.0

	// Ожидаем, что метод Transfer вызовет ошибку
	mockRepo.On("Transfer", mock.Anything, senderID, receiverID, amount).Return(errors.New("transfer failed"))

	err := service.Transfer(context.Background(), senderID, receiverID, amount)

	// Проверяем, что ошибка возвращена и она соответствует ожидаемой
	assert.Error(t, err)
	assert.Equal(t, "transfer failed", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetTransactions(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	// Входные данные
	userID := int64(1)
	expectedTransactions := []postgres.Transaction{
		{ID: 1, UserID: &userID, Amount: 100.0, TransactionType: "deposit"},
		{ID: 2, SenderID: &userID, ReceiverID: new(int64), Amount: 50.0, TransactionType: "transfer"},
	}

	// Ожидаем, что метод GetTransactions вернет список транзакций без ошибок
	mockRepo.On("GetTransactions", mock.Anything, userID).Return(expectedTransactions, nil)

	transactions, err := service.GetTransactions(context.Background(), userID)

	// Проверяем, что ошибок нет и транзакции совпадают с ожидаемыми
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, expectedTransactions, transactions)
	mockRepo.AssertExpectations(t)
}

func TestGetTransactions_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	userID := int64(1)

	// Ожидаем, что метод GetTransactions вызовет ошибку
	mockRepo.On("GetTransactions", mock.Anything, userID).Return([]postgres.Transaction{}, errors.New("db error"))

	transactions, err := service.GetTransactions(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	assert.Equal(t, []postgres.Transaction{}, transactions)
	mockRepo.AssertExpectations(t)
}
