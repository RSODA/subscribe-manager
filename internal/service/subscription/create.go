package subscription

import (
	"context"
	"strings"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/google/uuid"
)

func (s *subscriptionService) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	if sub.UserID == uuid.Nil || strings.TrimSpace(sub.ServiceName) == "" || sub.Price <= 0 || sub.StartDate.IsZero() {
		s.l.Errorw("invalid subscription data", "subscription", sub)
		return nil, apperrors.ErrInvalidSubscriptionData
	}

	res, err := s.repo.Create(ctx, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}
