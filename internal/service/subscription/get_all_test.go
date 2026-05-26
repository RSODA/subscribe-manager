package subscription

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/RSODA/subscribe-manager/internal/repository/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetAll(t *testing.T) {
	ctx := context.Background()

	mc := minimock.NewController(t)
	repoMock := mocks.NewSubscriptionRepositoryMock(mc)
	service := NewSubscriptionService(repoMock, zap.NewNop().Sugar())

	subs := []*domain.Subscription{
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			ServiceName: "Telegram Premium",
			Price:       399,
			StartDate:   time.Now().UTC(),
		},
		{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			ServiceName: "Yandex Plus",
			Price:       299,
			StartDate:   time.Now().UTC(),
		},
	}

	t.Run("successful get all", func(t *testing.T) {
		repoMock.GetAllMock.Return(subs, nil)

		res, err := service.GetAll(ctx)
		require.NoError(t, err)
		require.Equal(t, subs, res)
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("repository get all error")
		repoMock.GetAllMock.Return(nil, repoErr)

		res, err := service.GetAll(ctx)
		require.Nil(t, res)
		require.ErrorIs(t, err, repoErr)
	})
}
