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
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "wrong request")
		return
	}
	cafe := h.services.CafeList.GetCafeByID(cafeId)

	positions := h.services.CafeList.GetMenuByCafeID(cafeId)
	if len(positions) < 1 {
		newErrorResponse(c, http.StatusInternalServerError, "menu is empty")
		return
	}

	categories := []serv.Category{}
	name := h.services.CafeList.GetCategoriesByCafeID(cafeId)
	for _, category := range name {
		var dish_List []serv.Menu
		for _, pos := range positions {
			if pos.Category == category {
				dish_List = append(dish_List, pos)
			}
		}
		categories = append(categories, serv.Category{
			Category_Name: category,
			Menu_List: dish_List,
		})
	}

	nameAndMenu := &serv.CafeAndMenu{
		Cafe_Name: cafe.Name,
		Categories: categories,
	}
	c.JSON(http.StatusOK, nameAndMenu)
}
