package subscription

import (
	"bytes"
	"encoding/json"
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

func TestUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mc := minimock.NewController(t)
	serviceMock := mocks.NewSubscriptionServiceMock(mc)
	handler := NewHandler(serviceMock, zap.NewNop().Sugar())

	router := gin.New()
	handler.RegisterRoutes(router.Group(""))

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/subscriptions/not-uuid", bytes.NewBufferString(`{"price":100}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		id := uuid.New()
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/subscriptions/"+id.String(), bytes.NewBufferString("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("dto validation error", func(t *testing.T) {
		id := uuid.New()
		body := map[string]any{"price": 0}
		b, _ := json.Marshal(body)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/subscriptions/"+id.String(), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("success", func(t *testing.T) {
		id := uuid.New()
		userID := uuid.New()
		price := 300
		name := "Yandex Plus"
		body := map[string]any{"user_id": userID.String(), "service_name": name, "price": price, "start_date": "07-2026"}
		b, _ := json.Marshal(body)

		sub := &domain.Subscription{ID: id, UserID: userID, ServiceName: name, Price: int64(price), StartDate: time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)}
		serviceMock.UpdateMock.Return(sub, nil)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/subscriptions/"+id.String(), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service invalid data", func(t *testing.T) {
		id := uuid.New()
		body := map[string]any{"price": 100}
		b, _ := json.Marshal(body)
		serviceMock.UpdateMock.Return(nil, apperrors.ErrInvalidSubscriptionData)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/subscriptions/"+id.String(), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service internal error", func(t *testing.T) {
		id := uuid.New()
		body := map[string]any{"price": 100}
		b, _ := json.Marshal(body)
		serviceMock.UpdateMock.Return(nil, errors.New("boom"))
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/subscriptions/"+id.String(), bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
