package handler

import (
	"github.com/RSODA/subscribe-manager/internal/handler/middleware"
	subscriptionhandler "github.com/RSODA/subscribe-manager/internal/handler/subscription"
	"github.com/RSODA/subscribe-manager/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "github.com/RSODA/subscribe-manager/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	subscriptionService service.SubscriptionService
	l                   *zap.SugaredLogger
}

func newRouter(subscriptionService service.SubscriptionService, l *zap.SugaredLogger) *Router {
	return &Router{
		subscriptionService: subscriptionService,
		l:                   l,
	}
}

func (r *Router) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), middleware.ZapLogger(r.l))

	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/swagger/index.html")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	subscriptionHandler := subscriptionhandler.NewHandler(r.subscriptionService, r.l)
	subscriptionHandler.RegisterRoutes(router.Group(""))

	return router
}
