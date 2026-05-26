package subscription

import (
	"context"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/google/uuid"
)

func (s *subscriptionService) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	if sub.UserID == uuid.Nil {
		s.l.Errorw("invalid user ID", "userID", sub.UserID)
		return nil, apperrors.ErrInvalidSubscriptionData
	}

	if len(sub.ServiceName) == 0 || sub.Price <= 0 || sub.StartDate.IsZero() {
		s.l.Errorw("invalid subscription data", "subscription", sub)
		return nil, apperrors.ErrInvalidSubscriptionData
	}

	res, err := s.repo.Create(ctx, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}
