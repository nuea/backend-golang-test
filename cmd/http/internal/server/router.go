package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nuea/backend-golang-test/cmd/http/internal/handler"
	"github.com/nuea/backend-golang-test/internal/middleware"

	_ "github.com/nuea/backend-golang-test/cmd/http/internal/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerRouter(gin *gin.Engine, h *handler.Handlers, m *middleware.Middleware) {
	gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router := *gin.Group("/api/v1")
	{
		router.POST("/login", h.AuthHandler.Login)
		router.POST("/users", h.UserHandler.CreateUser)

		router.Use(m.Auth.Middleware())
		router.GET("/users", h.UserHandler.GetUsers)
		router.GET("/users/:id", h.UserHandler.GetUser)
		router.PATCH("/users/:id", h.UserHandler.UpdateUser)
		router.DELETE("/users/:id", h.UserHandler.DeleteUser)
	}
}
