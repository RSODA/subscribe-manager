package subscription

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/domain"
	"github.com/RSODA/subscribe-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mc := minimock.NewController(t)
	serviceMock := mocks.NewSubscriptionServiceMock(mc)
	handler := NewHandler(serviceMock, zap.NewNop().Sugar())

	router := gin.New()
	handler.RegisterRoutes(router.Group(""))

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/not-uuid", nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		sub := &domain.Subscription{ID: id, UserID: uuid.New(), ServiceName: "Telegram Premium", Price: 100, StartDate: time.Now().UTC()}
		serviceMock.GetByIDMock.Return(sub, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/"+id.String(), nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		id := uuid.New()
		serviceMock.GetByIDMock.Return(nil, apperrors.ErrSubscriptionNotFound)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/"+id.String(), nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		id := uuid.New()
		serviceMock.GetByIDMock.Return(nil, errors.New("boom"))

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/subscriptions/"+id.String(), nil)
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
