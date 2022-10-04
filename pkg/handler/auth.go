package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/callmehorhe/backtest/pkg/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid input body")
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
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
	user, err := h.services.GetUser(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
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
		"userId": user.Id_User,
		"phone":  user.Phone,
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

func (h *Handler) userProfile(c *gin.Context) {
	user_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "incorrect user id")
		return
	}
	user, err := h.services.Authorization.GetUserByID(user_id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "such user not found")
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) confirm(c *gin.Context) {
	code := c.Param("code")
	err := h.services.Authorization.ConfirmUser(code)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "wrong activation")
		return
	}
	c.Status(http.StatusOK)
}

type dataForgetPass struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (h *Handler) forgetPass(c *gin.Context) {
	var input dataForgetPass
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid input")
		return
	}
	if err := h.services.Authorization.ForgetPassword(input.Email, input.Phone); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) resetPassword(c *gin.Context) {
	type dataResetPass struct {
		Auth        string `json:"auth"`
		NewPassword string `json:"password"`
	}
	var input dataResetPass
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "invalid input")
		return
	}
	if err := h.services.Authorization.ResetPassword(input.Auth, input.NewPassword); err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
