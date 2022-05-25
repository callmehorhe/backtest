package handler

import (
	"net/http"
	"strconv"

	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/gin-gonic/gin"
	_ "gorm.io/driver/postgres"
)

func (h *Handler) orderSend(c *gin.Context) {
	var input models.Order
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}
	h.services.TGBot.SendOrder(input)

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) getOrderList(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "wrong request")
		return
	}
	orders := h.services.Order.GetOrdersByUser(user_id)
	c.JSON(http.StatusOK, orders)
}
