package handler

import (
	"net/http"
	"strconv"

	serv "github.com/callmehorhe/backtest"
	"github.com/gin-gonic/gin"
)

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
	cafe := h.services.CafeList.GetCafeByID(cafeId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "wrong request")
		return
	}

	positions := h.services.CafeList.GetMenuByCafeID(cafeId)
	if len(positions) < 1 {
		newErrorResponse(c, http.StatusInternalServerError, "menu is empty")
		return
	}
	nameAndMenu := &serv.CafeAndMenu{
		CafeName: cafe.Name,
		Menu:     positions,
	}
	c.JSON(http.StatusOK, nameAndMenu)
}
