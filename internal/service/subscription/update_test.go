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

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	mc := minimock.NewController(t)
	repoMock := mocks.NewSubscriptionRepositoryMock(mc)
	service := NewSubscriptionService(repoMock, zap.NewNop().Sugar())

	validSub := &domain.Subscription{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		ServiceName: "Telegram Premium",
		Price:       499,
		StartDate:   time.Now().UTC(),
	}

	updatedSub := &domain.Subscription{
		ID:          validSub.ID,
		UserID:      validSub.UserID,
		ServiceName: "Yandex Plus",
		Price:       299,
		StartDate:   validSub.StartDate,
	}

	t.Run("successful update", func(t *testing.T) {
		repoMock.UpdateMock.Return(updatedSub, nil)

		res, err := service.Update(ctx, validSub)
		require.NoError(t, err)
		require.Equal(t, updatedSub, res)
	})

	t.Run("invalid subscription ID", func(t *testing.T) {
		invalidSub := &domain.Subscription{ID: uuid.Nil}

		res, err := service.Update(ctx, invalidSub)
		require.Nil(t, res)
		require.ErrorIs(t, err, apperrors.ErrInvalidSubscriptionData)
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("repository update error")
		repoMock.UpdateMock.Return(nil, repoErr)

		res, err := service.Update(ctx, validSub)
		require.Nil(t, res)
		require.ErrorIs(t, err, repoErr)
	})
}
