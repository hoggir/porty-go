package routes

import (
	"net/http"
	"porty-go/controllers"
	"porty-go/repositories"
	"porty-go/services"

	"github.com/gin-gonic/gin"
)

// CharacterRoutes defines the character-related routes
func CharacterRoutes(r *gin.Engine) {
	repo, err := repositories.NewCharacterRepository()
	if err != nil {
		// handle the error appropriately
		r.GET("/characters", func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		})
		return
	}
	characterController := controllers.NewCharacterController(services.NewCharacterService(repo))

	r.GET("/characters", characterController.GetAllCharacters)
}
