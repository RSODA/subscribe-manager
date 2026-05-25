package subscription

import (
	"github.com/RSODA/subscribe-manager/internal/repository"
	"github.com/RSODA/subscribe-manager/internal/service"
	"go.uber.org/zap"
)

type subscriptionService struct {
	repo repository.SubscriptionRepository
	l    *zap.SugaredLogger
}

func NewSubscriptionService(repo repository.SubscriptionRepository, l *zap.SugaredLogger) service.SubscriptionService {
	return &subscriptionService{
		repo: repo,
		l:    l,
	}
}
