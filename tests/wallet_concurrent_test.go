package tests

import (
	"context"
	"database/sql"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shagabiev/wallet-service/internal/service"
	"github.com/shagabiev/wallet-service/internal/storage"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := storage.NewPostgres(
		"postgres://wallet:wallet@localhost:55432/wallet?sslmode=disable",
	)
	require.NoError(t, err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS wallets (id UUID PRIMARY KEY, balance BIGINT NOT NULL DEFAULT 0)`)
	require.NoError(t, err)

	_, err = db.Exec(`TRUNCATE TABLE wallets`)
	require.NoError(t, err)

	return db
}

func TestConcurrentDeposits(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := service.New(db)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	walletID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	_, err := db.Exec(`INSERT INTO wallets(id, balance) VALUES($1, 0)`, walletID)
	require.NoError(t, err)

	wg := sync.WaitGroup{}
	const requests = 1000

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = svc.UpdateBalance(ctx, walletID, "DEPOSIT", 1)
		}()
	}

	wg.Wait()

	balance, err := svc.GetBalance(ctx, walletID)
	require.NoError(t, err)
	require.Equal(t, int64(requests), balance)
}
