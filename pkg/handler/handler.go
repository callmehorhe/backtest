package handler

import (
	"github.com/callmehorhe/backtest/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	rot := gin.New()
	rot.Use()
	router := gin.Default()
	router.Use(CORSMiddleware)
	router.Static("/images", "./images")
	slash := router.Group("/")
	{
		auth := slash.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.POST("/sign-in", h.signIn)
			auth.GET("/sign-out", h.signOut)
			auth.POST("/sign-in-cafe", h.signInCafe)
			auth.GET("/confirm/:code", h.confirm)
			auth.GET("/forget-pass", h.forgetPass)
			auth.GET("/reset-pass", h.resetPassword)
		}
		api := slash.Group("/api", h.Auth)
		{
			cafes := api.Group("/cafes")
			{
				cafes.GET("/", h.getCafeList)
				cafes.GET("/:id", h.getMenuByCafeID)
				cafes.POST("/admin", h.changeMenu)
			}
			api.GET("/user/:id", h.userProfile)
			api.POST("/order", h.orderSend)
			api.GET("/orders/:id/:page", h.getOrderList)
		}
	}

	return router
}
