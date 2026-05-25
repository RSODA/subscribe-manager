package subscription

import (
	"context"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
)

func (s *subscriptionService) Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	if len(sub.ID) == 0 {
		s.l.Errorw("invalid subscription ID", "id", sub.ID)
		return nil, apperrors.ErrInvalidSubscriptionData
	}

	res, err := s.repo.Update(ctx, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}
