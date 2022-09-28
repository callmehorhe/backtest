package handler

import (
	"io"
	"os"

	"github.com/callmehorhe/backtest/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	if err := os.MkdirAll("./logs/gin", 0777); err != nil {
		logrus.Fatalf("Can't create log path: %v", err)
	}
	logFile, err := os.OpenFile("./logs/gin/logs.txt", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrus.Fatalf("cant create log file: %v", err)
	}
	defer logFile.Close()
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(logFile, os.Stderr)

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
