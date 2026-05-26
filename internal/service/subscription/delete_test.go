package subscription

import (
	"context"
	"errors"
	"testing"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/repository/mocks"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDelete(t *testing.T) {
	mc := minimock.NewController(t)
	repoMock := mocks.NewSubscriptionRepositoryMock(mc)
	service := NewSubscriptionService(repoMock, zap.NewNop().Sugar())

	t.Run("successful delete", func(t *testing.T) {
		repoMock.DeleteMock.Return(nil)
		err := service.Delete(context.Background(), uuid.New().String())
		require.NoError(t, err)
	})

	t.Run("invalid subscription ID", func(t *testing.T) {
		err := service.Delete(context.Background(), "invalid-uuid")
		require.ErrorIs(t, err, apperrors.ErrInvalidSubscriptionData)
	})

	t.Run("repository error", func(t *testing.T) {
		repoMock.DeleteMock.Return(errors.New("delete error"))
		err := service.Delete(context.Background(), uuid.New().String())
		require.Error(t, err)
	})
}
