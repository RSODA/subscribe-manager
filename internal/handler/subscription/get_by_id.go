package subscription

import (
	"net/http"

	"github.com/RSODA/subscribe-manager/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Получить подписку по ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "UUID подписки"
// @Success      200  {object}  dto.SubscriptionResponse
// @Failure      400  {object}  dto.BadRequestErrorResponse
// @Failure      404  {object}  dto.NotFoundErrorResponse
// @Failure      500  {object}  dto.InternalErrorResponse
// @Router       /subscriptions/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		badRequest(c, "invalid subscription id")
		return
	}

	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NewSubscriptionResponse(res))
}
