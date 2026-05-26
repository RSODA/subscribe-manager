package subscription

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/RSODA/subscribe-manager/internal/repository/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	mc := minimock.NewController(t)
	repoMock := mocks.NewSubscriptionRepositoryMock(mc)

	service := NewSubscriptionService(repoMock, zap.NewNop().Sugar())

	validSub := &domain.Subscription{
		UserID:      uuid.New(),
		ServiceName: "Telegram Premium",
		Price:       100,
		StartDate:   time.Now(),
	}

	t.Run("valid subscription", func(t *testing.T) {
		repoMock.CreateMock.Return(validSub, nil)
		res, err := service.Create(ctx, validSub)
		require.Equal(t, validSub, res)
		require.NoError(t, err)
	})

	t.Run("invalid data subscription name", func(t *testing.T) {
		invalidSub := &domain.Subscription{
			UserID:      uuid.New(),
			ServiceName: "",
			Price:       100,
			StartDate:   time.Now(),
		}
		res, err := service.Create(ctx, invalidSub)
		require.Nil(t, res)
		require.ErrorIs(t, err, apperrors.ErrInvalidSubscriptionData)
	})

	t.Run("invalid data subscription price", func(t *testing.T) {
		invalidSub := &domain.Subscription{
			UserID:      uuid.New(),
			ServiceName: "Telegram Premium",
			Price:       -10,
			StartDate:   time.Now(),
		}
		res, err := service.Create(ctx, invalidSub)
		require.Nil(t, res)
		require.ErrorIs(t, err, apperrors.ErrInvalidSubscriptionData)
	})

	t.Run("UUID nil", func(t *testing.T) {
		invalidSub := &domain.Subscription{
			UserID:      uuid.Nil,
			ServiceName: "Telegram Premium",
			Price:       100,
			StartDate:   time.Now(),
		}
		res, err := service.Create(ctx, invalidSub)
		require.Nil(t, res)
		require.ErrorIs(t, err, apperrors.ErrInvalidSubscriptionData)
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("repository create error")
		repoMock.CreateMock.Return(nil, repoErr)

		res, err := service.Create(ctx, validSub)
		require.Nil(t, res)
		require.ErrorIs(t, err, repoErr)
	})
}
