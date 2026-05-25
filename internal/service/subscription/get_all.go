package subscription

import (
	"context"

	"github.com/RSODA/subscribe-manager/internal/domain"
)

func (s *subscriptionService) GetAll(ctx context.Context) ([]*domain.Subscription, error) {
	res, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
