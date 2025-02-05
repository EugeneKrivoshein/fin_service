package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Deposit(ctx context.Context, userID int64, amount float64) error
	Transfer(ctx context.Context, senderID, receiverID int64, amount float64) error
	GetTransactions(ctx context.Context, userID int64) ([]Transaction, error)
}

type Transaction struct {
	ID              int64     `json:"id"`
	UserID          *int64    `json:"user_id,omitempty"`
	SenderID        *int64    `json:"sender_id,omitempty"`
	ReceiverID      *int64    `json:"receiver_id,omitempty"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}

type RepositoryImpl struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *RepositoryImpl {
	return &RepositoryImpl{pool: pool}
}

// Пополняет баланс пользователя и создаёт транзакцию типа "deposit"
func (r *RepositoryImpl) Deposit(ctx context.Context, userID int64, amount float64) (err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	// Откат транзакции при ошибке
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	updateQuery := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	ct, err := tx.Exec(ctx, updateQuery, amount, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	insertQuery := `
		INSERT INTO transactions (user_id, amount, transaction_type)
		VALUES ($1, $2, 'deposit')
	`
	_, err = tx.Exec(ctx, insertQuery, userID, amount)
	if err != nil {
		return fmt.Errorf("failed to insert deposit transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit deposit transaction: %w", err)
	}

	return nil
}

// Переводит деньги от одного пользователя к другому
func (r *RepositoryImpl) Transfer(ctx context.Context, senderID, receiverID int64, amount float64) (err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Проверяем, что у отправителя достаточно средств
	var senderBalance float64
	err = tx.QueryRow(ctx, `SELECT balance FROM users WHERE id = $1`, senderID).Scan(&senderBalance)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}
	if senderBalance < amount {
		return errors.New("insufficient funds")
	}

	updateSenderQuery := `UPDATE users SET balance = balance - $1 WHERE id = $2`
	ct, err := tx.Exec(ctx, updateSenderQuery, amount, senderID)
	if err != nil {
		return fmt.Errorf("failed to update sender balance: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errors.New("sender not found")
	}

	updateReceiverQuery := `UPDATE users SET balance = balance + $1 WHERE id = $2`
	ct, err = tx.Exec(ctx, updateReceiverQuery, amount, receiverID)
	if err != nil {
		return fmt.Errorf("failed to update receiver balance: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return errors.New("receiver not found")
	}

	insertQuery := `
		INSERT INTO transactions (user_id, sender_id, receiver_id, amount, transaction_type)
		VALUES ($1, $2, $3, $4, 'transfer')
	`
	_, err = tx.Exec(ctx, insertQuery, senderID, senderID, receiverID, amount)
	if err != nil {
		return fmt.Errorf("failed to insert transfer transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transfer transaction: %w", err)
	}

	return nil
}

// Получает 10 последних транзакций для указанного пользователя
func (r *RepositoryImpl) GetTransactions(ctx context.Context, userID int64) ([]Transaction, error) {
	query := `
		SELECT id, user_id, sender_id, receiver_id, amount, transaction_type, created_at
		FROM transactions
		WHERE user_id = $1 OR sender_id = $1 OR receiver_id = $1
		ORDER BY created_at DESC
		LIMIT 10
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.SenderID, &t.ReceiverID, &t.Amount, &t.TransactionType, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
