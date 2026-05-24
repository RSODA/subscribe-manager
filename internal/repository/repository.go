package repository

import (
	"context"
	"time"

	"github.com/RSODA/subscribe-manager/internal/dto"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *dto.Subscription) (*dto.Subscription, error)
	GetByID(ctx context.Context, id string) (*dto.Subscription, error)
	GetAll(ctx context.Context) ([]*dto.Subscription, error)
	Update(ctx context.Context, sub *dto.Subscription) (*dto.Subscription, error)
	Delete(ctx context.Context, id string) error
	TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, from *time.Time, to *time.Time) (int, error)
}
