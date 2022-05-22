package handler

import (
	"net/http"
	"strconv"
	"time"

	serv "github.com/callmehorhe/backtest"
	"github.com/gin-gonic/gin"
)

type cafeData struct {
	ID       int
	Password string
}

func (h *Handler) signInCafe(c *gin.Context) {
	var input cafeData

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid input")
		return
	}
	cafe, err := h.services.GetCafe(input.ID, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := h.services.CafeList.CafeGenerateToken(input.ID, input.Password)
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
		"cafeId": cafe.Chat_ID,
	})
}

func (h *Handler) getCafeList(c *gin.Context) {
	cafes := h.services.CafeList.GetCafeList()
	if len(cafes) < 1 {
		newErrorResponse(c, http.StatusInternalServerError, "cafes not found")
		return
	}
	c.JSON(http.StatusOK, cafes)
}

func (h *Handler) getMenuByCafeID(c *gin.Context) {
	cafeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "wrong request")
		return
	}
	cafe := h.services.CafeList.GetCafeByID(cafeId)

	positions := h.services.CafeList.GetMenuByCafeID(cafeId)

	var dish_List []serv.Menu
	categories := h.services.CafeList.GetCategoriesByCafeID(cafeId)
	for _, category := range categories {
		for _, pos := range positions {
			if pos.Category == category {
				dish_List = append(dish_List, pos)
			}
		}
	}
	if len(dish_List) < 1 {
		newErrorResponse(c, http.StatusInternalServerError, "menu is empty")
		return
	}

	nameAndMenu := &serv.CafeAndMenu{
		Cafe_Name:  cafe.Name,
		Categories: categories,
		Menu:       dish_List,
	}
	c.JSON(http.StatusOK, nameAndMenu)
}
