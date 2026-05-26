package subscription

import (
	"context"
	"testing"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/RSODA/subscribe-manager/internal/repository/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	mc := minimock.NewController(t)
	repoMock := mocks.NewSubscriptionRepositoryMock(mc)
	service := NewSubscriptionService(repoMock, zap.NewNop().Sugar())

	validID := uuid.New().String()
	expectedSub := &domain.Subscription{
		ID:          uuid.MustParse(validID),
		ServiceName: "Yandex Plus",
		Price:       100,
	}

	t.Run("valid subscription ID", func(t *testing.T) {
		repoMock.GetByIDMock.Return(expectedSub, nil)
		res, err := service.GetByID(ctx, validID)
		require.Equal(t, expectedSub, res)
		require.NoError(t, err)
	})

	t.Run("invalid subscription ID", func(t *testing.T) {
		res, err := service.GetByID(ctx, "invalid-uuid")
		require.Nil(t, res)
		require.ErrorIs(t, err, apperrors.ErrInvalidSubscriptionData)
	})

	t.Run("not found", func(t *testing.T) {
		repoMock.GetByIDMock.Return(nil, apperrors.ErrSubscriptionNotFound)
		res, err := service.GetByID(ctx, validID)
		require.Nil(t, res)
		require.ErrorIs(t, err, apperrors.ErrSubscriptionNotFound)
	})
}
