package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) UpdateBalance(ctx context.Context, walletID uuid.UUID, operation string, amount int64) error {
	delta := amount
	if operation == "WITHDRAW" {
		delta = -amount
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var balance int64
	if err = tx.QueryRowContext(ctx,
		`SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`, walletID).
		Scan(&balance); err != nil {
		tx.Rollback()
		return err
	}

	if balance+delta < 0 {
		tx.Rollback()
		return ErrInsufficientFunds
	}

	if _, err = tx.ExecContext(ctx,
		`UPDATE wallets SET balance = balance + $1 WHERE id = $2`, delta, walletID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Service) GetBalance(ctx context.Context, walletID uuid.UUID) (int64, error) {
	var balance int64
	err := s.db.QueryRowContext(
		ctx,
		`SELECT balance FROM wallets WHERE id = $1`,
		walletID,
	).Scan(&balance)

	return balance, err
}
