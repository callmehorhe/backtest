package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/gin-gonic/gin"
)

type CafeData struct {
	CafeId   string `json:"cafeId"`
	Password string `json:"password"`
}

func (h *Handler) signInCafe(c *gin.Context) {
	var input CafeData

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid input")
		return
	}
	id, err := strconv.Atoi(input.CafeId)
	if err != nil {
		log.Println(err)
	}

	cafe, err := h.services.GetCafe(id, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := h.services.CafeList.CafeGenerateToken(id, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Writer.Header().Set("Authorization", token)
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, map[string]interface{}{
		"token":  token,
		"cafeId": cafe.Id_Cafe,
	})
}

type CafeProperties struct {
	Cafe models.Cafe   `json:"cafe"`
	Menu []models.Menu `json:"menu"`
}

func (h *Handler) changeMenu(c *gin.Context) {
	var input CafeProperties
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid input")
	}
	err := h.services.CafeList.UpdateCafe(input.Cafe)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid write to db")
	}
	h.services.CafeList.UpdateMenu(input.Menu, input.Cafe.Name)
	c.Status(http.StatusOK)
}
