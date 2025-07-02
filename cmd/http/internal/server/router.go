package server

import (
	"github.com/gin-gonic/gin"
	"github.com/nuea/backend-golang-test/cmd/http/internal/handler"
)

func registerRouter(gin *gin.Engine, h *handler.Handlers) {
	router := *gin.Group("/api/v1")
	{
		router.POST("/register", h.UserHandler.Register)
		router.GET("/users", h.UserHandler.GetUsers)
		router.GET("/users/:id", h.UserHandler.GetUser)
		router.PATCH("/users/:id", h.UserHandler.UpdateUser)
		router.DELETE("/users/:id", h.UserHandler.DeleteUser)
	}
}
