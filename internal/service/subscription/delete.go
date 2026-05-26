package subscription

import (
	"context"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/google/uuid"
)

func (s *subscriptionService) Delete(ctx context.Context, id string) error {
	if _, err := uuid.Parse(id); err != nil {
		s.l.Errorw("invalid subscription ID format", "id", id)
		return apperrors.ErrInvalidSubscriptionData
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
