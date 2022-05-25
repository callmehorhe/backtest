package handler

import (
	"net/http"
	"strconv"

	"github.com/callmehorhe/backtest/pkg/models"
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
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "wrong request")
		return
	}
	cafe := h.services.CafeList.GetCafeByID(cafeId)

	positions := h.services.CafeList.GetMenuByCafeID(cafeId)

	var dish_List []models.Menu
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

	nameAndMenu := &models.CafeAndMenu{
		Cafe_Name:  cafe.Name,
		Categories: categories,
		Menu:       dish_List,
	}
	c.JSON(http.StatusOK, nameAndMenu)
}
