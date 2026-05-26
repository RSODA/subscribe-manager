package subscription

import (
	"net/http"

	"github.com/RSODA/subscribe-manager/internal/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := req.ToDomain()
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	res, err := h.service.Create(c.Request.Context(), sub)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewSubscriptionResponse(res))
}
