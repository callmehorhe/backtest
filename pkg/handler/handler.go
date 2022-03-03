package handler

import (
	"github.com/callmehorhe/backtest/pkg/service"
	"github.com/gin-contrib/cors"
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
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	slash := router.Group("/")
	{
		api := slash.Group("/api")
		{
			auth := api.Group("/auth")
			{
				auth.POST("/sign-up", h.signUp)
				auth.POST("/sign-in", h.signIn)
				auth.GET("/sign-out", h.signOut)

			}
			cafes := api.Group("/cafes")
			{
				cafes.GET("/", h.getCafeList)
				cafes.GET("/:id", h.getMenuByCafeID)
			}
			api.POST("/order", h.orderSend)
			api.GET("/orders/:id", h.getOrderList)
		}
	}
	return router
}
