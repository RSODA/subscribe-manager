package subscription

import (
	"net/http"

	"github.com/RSODA/subscribe-manager/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Обновить подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      string                         true  "UUID подписки"
// @Param        subscription  body      dto.UpdateSubscriptionRequest  true  "Данные для обновления"
// @Success      200           {object}  dto.SubscriptionResponse
// @Failure      400           {object}  dto.BadRequestErrorResponse
// @Failure      404           {object}  dto.NotFoundErrorResponse
// @Failure      500           {object}  dto.InternalErrorResponse
// @Router       /subscriptions/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		badRequest(c, "invalid subscription id")
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := req.ToDomain(parsedID)
	if err != nil {
		badRequest(c, err.Error())
		return
	}

	res, err := h.service.Update(c.Request.Context(), sub)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NewSubscriptionResponse(res))
}
