package subscription

import (
	"context"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
)

func (s *subscriptionService) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	if len(sub.UserID) == 0 || len(sub.ServiceName) == 0 || sub.Price <= 0 || sub.StartDate.IsZero() {
		s.l.Errorw("invalid subscription data", "subscription", sub)
		return nil, apperrors.ErrInvalidSubscriptionData
	}

	res, err := s.repo.Create(ctx, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}
