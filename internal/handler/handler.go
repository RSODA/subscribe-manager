package handler

import (
	"github.com/RSODA/subscribe-manager/internal/service"
	"go.uber.org/zap"
)

func NewRouter(subscriptionService service.SubscriptionService, l *zap.SugaredLogger) *Router {
	return newRouter(subscriptionService, l)
}
