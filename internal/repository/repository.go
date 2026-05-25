package repository

import (
	"context"
	"time"

	"github.com/RSODA/subscribe-manager/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	GetByID(ctx context.Context, id string) (*domain.Subscription, error)
	GetAll(ctx context.Context) ([]*domain.Subscription, error)
	Update(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	Delete(ctx context.Context, id string) error
	TotalCost(ctx context.Context, userID *string, serviceName *string, from *time.Time, to *time.Time) (int, error)
}
