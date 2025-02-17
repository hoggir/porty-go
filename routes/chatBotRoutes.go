package routes

import (
	"porty-go/controllers"
	middleware "porty-go/middlewares"

	"github.com/gin-gonic/gin"
)

// CharacterRoutes defines the character-related routes
func AiRoutes(r *gin.Engine) {
	protected := r.Group("/chat")
	protected.Use(middleware.JWTAuth())
	{
		protected.POST("/", controllers.ChatAi)
	}
}
