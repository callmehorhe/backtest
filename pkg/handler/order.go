package handler

import (
	serv "github.com/callmehorhe/backtest"
	"github.com/gin-gonic/gin"
)

type orderData struct {
	Cafe string `json:"cafe"`
	Positions []serv.Position
	Address string `json:"address"`
}

func (h *Handler) OrderSending(c *gin.Context) {
	
}

