package subscription

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RSODA/subscribe-manager/internal/repository/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTotalCost(t *testing.T) {
	ctx := context.Background()

	mc := minimock.NewController(t)
	repoMock := mocks.NewSubscriptionRepositoryMock(mc)
	service := NewSubscriptionService(repoMock, zap.NewNop().Sugar())

	userID := uuid.New().String()
	serviceName := "Telegram Premium"
	from := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, time.December, 1, 0, 0, 0, 0, time.UTC)

	t.Run("successful total cost", func(t *testing.T) {
		repoMock.TotalCostMock.Return(698, nil)

		res, err := service.TotalCost(ctx, &userID, &serviceName, &from, &to)
		require.NoError(t, err)
		require.Equal(t, 698, res)
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("repository total cost error")
		repoMock.TotalCostMock.Return(0, repoErr)

		res, err := service.TotalCost(ctx, &userID, &serviceName, &from, &to)
		require.Zero(t, res)
		require.ErrorIs(t, err, repoErr)
	})
}
