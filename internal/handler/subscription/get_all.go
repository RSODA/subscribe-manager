package subscription

import (
	"net/http"

	"github.com/RSODA/subscribe-manager/internal/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAll(c *gin.Context) {
	res, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NewSubscriptionResponses(res))
}
