package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/shagabiev/wallet-service/internal/service"
	"github.com/shagabiev/wallet-service/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestConcurrentDeposits(t *testing.T) {
	db, err := storage.NewPostgres(
		"postgres://wallet:wallet@localhost:5432/wallet?sslmode=disable",
	)
	require.NoError(t, err)

	svc := service.New(db)
	ctx := context.Background()

	walletID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	wg := sync.WaitGroup{}
	requests := 1000

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
