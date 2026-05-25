package subscription

import (
	"context"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
)

func (s *subscriptionService) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	if len(id) == 0 {
		s.l.Errorw("invalid subscription ID", "id", id)
		return nil, apperrors.ErrInvalidSubscriptionData
	}

	res, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}
