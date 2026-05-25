package subscription

import (
	"context"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
)

func (s *subscriptionService) Delete(ctx context.Context, id string) error {
	if len(id) == 0 {
		s.l.Errorw("invalid subscription ID", "id", id)
		return apperrors.ErrInvalidSubscriptionData
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
