package subscription

import (
	"context"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/google/uuid"
)

func (s *subscriptionService) GetByID(ctx context.Context, id string) (*domain.Subscription, error) {
	if _, err := uuid.Parse(id); err != nil {
		s.l.Errorw("invalid subscription ID format", "id", id)
		return nil, apperrors.ErrInvalidSubscriptionData
	}
	res, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}
