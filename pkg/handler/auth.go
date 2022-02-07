package handler

import (
	"net/http"

	serv "github.com/callmehorhe/backtest"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input serv.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid input body")
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type data struct {
	email    string `json:"email"`
	password string `json:"password"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input data

	if err := c.BindJSON(&input); err != nil {
		
	}
}
