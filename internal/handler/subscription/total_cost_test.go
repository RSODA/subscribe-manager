package subscription

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RSODA/subscribe-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTotalCost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mc := minimock.NewController(t)
	serviceMock := mocks.NewSubscriptionServiceMock(mc)
	handler := NewHandler(serviceMock, zap.NewNop().Sugar())
	router := gin.New()
	handler.RegisterRoutes(router.Group(""))

	t.Run("invalid user id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/total?user_id=bad", nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid from date", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/total?from=13-2026", nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid to date", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/total?to=bad", nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("success without filters", func(t *testing.T) {
		serviceMock.TotalCostMock.Return(1234, nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/total", nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("success with all filters", func(t *testing.T) {
		serviceMock.TotalCostMock.Return(698, nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/total?user_id=0e7a1e5f-94f0-4968-838a-e0329b0d556e&service_name=Telegram%20Premium&from=01-2026&to=12-2026", nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		serviceMock.TotalCostMock.Return(0, errors.New("boom"))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/total", nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
