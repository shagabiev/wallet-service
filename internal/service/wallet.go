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

func (s *Service) UpdateBalance(
	ctx context.Context,
	walletID uuid.UUID,
	operation string,
	amount int64,
) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance int64
	err = tx.QueryRowContext(
		ctx,
		`SELECT balance FROM wallets WHERE id = $1 FOR UPDATE`,
		walletID,
	).Scan(&balance)
	if err != nil {
		return err
	}

	delta := amount
	if operation == "WITHDRAW" {
		delta = -amount
	}

	if balance+delta < 0 {
		return ErrInsufficientFunds
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE wallets SET balance = balance + $1 WHERE id = $2`,
		delta,
		walletID,
	)
	if err != nil {
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
