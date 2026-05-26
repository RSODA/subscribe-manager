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

func TestCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mc := minimock.NewController(t)
	serviceMock := mocks.NewSubscriptionServiceMock(mc)
	handler := NewHandler(serviceMock, zap.NewNop().Sugar())

	router := gin.New()
	handler.RegisterRoutes(router.Group(""))

	t.Run("successful create", func(t *testing.T) {
		expectedSub := &domain.Subscription{
			ID:          uuid.New(),
			UserID:      uuid.New(),
			ServiceName: "Telegram Premium",
			Price:       100,
			StartDate:   time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC),
		}
		serviceMock.CreateMock.Return(expectedSub, nil)

		body := map[string]interface{}{
			"user_id":      expectedSub.UserID.String(),
			"service_name": expectedSub.ServiceName,
			"price":        expectedSub.Price,
			"start_date":   expectedSub.StartDate.Format("01-2006"),
		}
		bodyBytes, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var res map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &res)
		require.NoError(t, err)
		require.Equal(t, "Telegram Premium", res["service_name"])
	})

	t.Run("invalid json", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBufferString("{invalid"))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("dto validation error", func(t *testing.T) {
		body := map[string]interface{}{
			"user_id":      uuid.New().String(),
			"service_name": "Telegram Premium",
			"price":        100,
			"start_date":   "bad-date",
		}
		bodyBytes, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service invalid data error", func(t *testing.T) {
		serviceMock.CreateMock.Return(nil, apperrors.ErrInvalidSubscriptionData)

		body := map[string]interface{}{
			"user_id":      uuid.New().String(),
			"service_name": "Yandex Plus",
			"price":        300,
			"start_date":   "07-2026",
		}
		bodyBytes, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service internal error", func(t *testing.T) {
		serviceMock.CreateMock.Return(nil, errors.New("boom"))

		body := map[string]interface{}{
			"user_id":      uuid.New().String(),
			"service_name": "Yandex Plus",
			"price":        300,
			"start_date":   "07-2026",
		}
		bodyBytes, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
