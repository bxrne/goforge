package routes

import (
	"github.com/bxrne/goforge/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(api *gin.RouterGroup) {
	// User routes
	users := api.Group("/users")
	{
		users.GET("", handlers.GetUsers)
		users.POST("", handlers.CreateUser)
	}
}
