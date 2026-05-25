package subscription

import (
	"context"
	"time"
)

func (s *subscriptionService) TotalCost(ctx context.Context, userID *string, serviceName *string, from *time.Time, to *time.Time) (int, error) {
	res, err := s.repo.TotalCost(ctx, userID, serviceName, from, to)
	if err != nil {
		return 0, err
	}
	return res, nil
}
