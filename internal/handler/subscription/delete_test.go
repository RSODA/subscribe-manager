package subscription

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mc := minimock.NewController(t)
	serviceMock := mocks.NewSubscriptionServiceMock(mc)
	handler := NewHandler(serviceMock, zap.NewNop().Sugar())
	router := gin.New()
	handler.RegisterRoutes(router.Group(""))

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/subscriptions/not-uuid", nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()
		serviceMock.DeleteMock.Return(nil)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/subscriptions/"+id, nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		id := uuid.New().String()
		serviceMock.DeleteMock.Return(apperrors.ErrSubscriptionNotFound)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/subscriptions/"+id, nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		id := uuid.New().String()
		serviceMock.DeleteMock.Return(errors.New("boom"))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/subscriptions/"+id, nil)
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
