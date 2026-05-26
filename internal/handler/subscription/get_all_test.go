package subscription

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/RSODA/subscribe-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mc := minimock.NewController(t)
	serviceMock := mocks.NewSubscriptionServiceMock(mc)
	handler := NewHandler(serviceMock, zap.NewNop().Sugar())

	router := gin.New()
	handler.RegisterRoutes(router.Group(""))

	t.Run("successful get all", func(t *testing.T) {
		subs := []*domain.Subscription{
			{ID: uuid.New(), UserID: uuid.New(), ServiceName: "Telegram Premium", Price: 399, StartDate: time.Now().UTC()},
			{ID: uuid.New(), UserID: uuid.New(), ServiceName: "Yandex Plus", Price: 299, StartDate: time.Now().UTC()},
		}
		serviceMock.GetAllMock.Return(subs, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		serviceMock.GetAllMock.Return(nil, errors.New("boom"))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
