package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetCafeList(c *gin.Context) {
	cafes := h.services.CafeList.GetCafeList()
	if len(cafes) < 1 {
		newErrorResponse(c, http.StatusInternalServerError, "cafes not found")
		return
	}
	c.JSON(http.StatusOK, cafes)
}