package subscription

import (
	"context"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/google/uuid"
)

func (s *subscriptionService) Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	if sub.ID == uuid.Nil {
		s.l.Errorw("invalid subscription ID", "id", sub.ID)
		return nil, apperrors.ErrInvalidSubscriptionData
	}

	res, err := s.repo.Update(ctx, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}
