package handler

import (
	"net/http"
	"time"

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
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input data

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid input")
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, cookie)

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) signOut(c *gin.Context) {
	cookie := &http.Cookie{
		Name:    "token",
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-time.Hour),
	}
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "logged out",
	})
}

func (h *Handler) auth(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "unauth user")
		return
	}
	id, err := h.services.ParseToken(token)
	if err != nil {
		h.signOut(c)
		newErrorResponse(c, http.StatusUnauthorized, "unauth user")
		return
	}
	user, err := h.services.GetUserByID(id)
	if err != nil {
		h.signOut(c)
		newErrorResponse(c, http.StatusUnauthorized, "user not found")
		return
	}
	c.JSON(http.StatusOK, user)
}
