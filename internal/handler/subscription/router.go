package subscription

import (
	"github.com/RSODA/subscribe-manager/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service service.SubscriptionService
	l       *zap.SugaredLogger
}

func NewHandler(subscriptionService service.SubscriptionService, l *zap.SugaredLogger) *Handler {
	return &Handler{
		service: subscriptionService,
		l:       l,
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	subscriptions := r.Group("/subscriptions")

	subscriptions.POST("", h.Create)
	subscriptions.GET("", h.GetAll)
	subscriptions.GET("/total", h.TotalCost)
	subscriptions.GET("/:id", h.GetByID)
	subscriptions.PUT("/:id", h.Update)
	subscriptions.DELETE("/:id", h.Delete)
}
