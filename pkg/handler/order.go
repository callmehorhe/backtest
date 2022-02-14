package handler

import (
	"net/http"

	serv "github.com/callmehorhe/backtest"
	"github.com/gin-gonic/gin"
)


func (h *Handler) OrderSend(c *gin.Context) {
	var input serv.Order

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request")
		return 
	}
	
	c.AbortWithStatus(http.StatusOK)
}

