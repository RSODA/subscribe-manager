package subscription

import (
	"net/http"
	"strings"
	"time"

	"github.com/RSODA/subscribe-manager/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) TotalCost(c *gin.Context) {
	userID, ok := optionalUUID(c, "user_id")
	if !ok {
		return
	}

	serviceName := optionalString(c.Query("service_name"))

	from, ok := optionalDate(c, "from")
	if !ok {
		return
	}
	to, ok := optionalDate(c, "to")
	if !ok {
		return
	}

	res, err := h.service.TotalCost(c.Request.Context(), userID, serviceName, from, to)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NewTotalCostResponse(res))
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func optionalUUID(c *gin.Context, name string) (*string, bool) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return nil, true
	}

	if _, err := uuid.Parse(value); err != nil {
		badRequest(c, name+" must be a valid uuid")
		return nil, false
	}

	return &value, true
}

func optionalDate(c *gin.Context, name string) (*time.Time, bool) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return nil, true
	}

	parsed, err := parseMonthYear(value)
	if err != nil {
		badRequest(c, name+" must use MM-YYYY format")
		return nil, false
	}

	return &parsed, true
}
