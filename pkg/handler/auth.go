package handler

import (
	"net/http"
	"strings"
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
	user := h.services.GetUser(input.Email, input.Password)
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
		"token":  token,
		"userId": user.Id_User,
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
	auth := c.GetHeader("Authorization")
	if len(auth) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "authorization header is not provided")
		c.AbortWithStatusJSON(http.StatusUnauthorized, "authorization header is not provided")
		return
	}

	fields := strings.Fields(auth)
	if len(fields) < 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
		c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid authorization header format")
		return
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != "bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "unsupported authorization type")
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unsupported authorization type")
		return
	}

	accessToken := fields[1]
	id, err := h.services.ParseToken(accessToken)
	if err != nil {
		h.signOut(c)
		newErrorResponse(c, http.StatusUnauthorized, "unauth user")
		c.AbortWithStatusJSON(http.StatusUnauthorized, "unauth user")
		return
	}

	user, err := h.services.GetUserByID(id)
	if err != nil {
		h.signOut(c)
		newErrorResponse(c, http.StatusUnauthorized, "user not found")
		return
	}
	c.Set("isLoggedIn", true)
	c.Set("userId", user.Id_User)
	c.Set("username", user.Name)
	c.Next()
}
