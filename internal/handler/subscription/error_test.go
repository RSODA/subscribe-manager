package subscription

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/RSODA/subscribe-manager/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHelperFunctions(t *testing.T) {
	t.Run("optional string", func(t *testing.T) {
		res := optionalString("   ")
		require.Nil(t, res)
		res = optionalString(" Telegram Premium ")
		require.NotNil(t, res)
		require.Equal(t, "Telegram Premium", *res)
	})

	t.Run("parse month year", func(t *testing.T) {
		parsed, err := parseMonthYear("07-2026")
		require.NoError(t, err)
		require.Equal(t, time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC), parsed)

		_, err = parseMonthYear("99-2026")
		require.Error(t, err)
	})
}

func TestWriteError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mc := minimock.NewController(t)
	serviceMock := mocks.NewSubscriptionServiceMock(mc)
	handler := NewHandler(serviceMock, zap.NewNop().Sugar())

	t.Run("bad request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.writeError(c, apperrors.ErrInvalidSubscriptionData)
		require.Equal(t, http.StatusBadRequest, w.Code)

		var body map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &body)
		require.NoError(t, err)
		require.Equal(t, apperrors.ErrInvalidSubscriptionData.Error(), body["error"])
	})

	t.Run("not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.writeError(c, apperrors.ErrSubscriptionNotFound)
		require.Equal(t, http.StatusNotFound, w.Code)

		var body map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &body)
		require.NoError(t, err)
		require.Equal(t, apperrors.ErrSubscriptionNotFound.Error(), body["message"])
	})

	t.Run("internal", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		handler.writeError(c, errors.New("boom"))
		require.Equal(t, http.StatusInternalServerError, w.Code)

		var body map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &body)
		require.NoError(t, err)
		require.Equal(t, "somthing internal error", body["message"])
	})

	t.Run("bad request helper", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		badRequest(c, "invalid")
		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}
