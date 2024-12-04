package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Register user routes
	UserRoutes(r)
	// Register character routes
	CharacterRoutes(r)
}
