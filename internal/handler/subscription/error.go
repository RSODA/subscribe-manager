package subscription

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/RSODA/subscribe-manager/internal/apperrors"
	"github.com/gin-gonic/gin"
)

const monthYearLayout = "01-2006"

func (h *Handler) writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrInvalidSubscriptionData):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, apperrors.ErrSubscriptionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		h.l.Errorw("subscription handler error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func badRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func parseMonthYear(value string) (time.Time, error) {
	parsed, err := time.Parse(monthYearLayout, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(parsed.Year(), parsed.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}
