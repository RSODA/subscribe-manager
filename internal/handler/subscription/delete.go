package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary      Удалить подписку
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "UUID подписки"
// @Success      204  "No Content"
// @Failure      400  {object}  dto.BadRequestErrorResponse
// @Failure      404  {object}  dto.NotFoundErrorResponse
// @Failure      500  {object}  dto.InternalErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		badRequest(c, "invalid subscription id")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.writeError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
