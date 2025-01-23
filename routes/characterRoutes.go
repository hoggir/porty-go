package routes

import (
	"fmt"
	"porty-go/controllers"
	"porty-go/repositories"
	"porty-go/services"

	"github.com/gin-gonic/gin"
)

// CharacterRoutes defines the character-related routes
func CharacterRoutes(r *gin.Engine) {
	repo, err := repositories.NewCharacterRepository()
	if err != nil {
		fmt.Println("Failed to create a new character repository: ", err)
		r.Use(func(c *gin.Context) {
			c.JSON(500, gin.H{
				"message": "Failed to create a new character repository: " + err.Error(),
				"status":  "error",
			})
			c.Abort()
		})
	}
	characterController := controllers.NewCharacterController(services.NewCharacterService(repo))

	r.GET("/characters", characterController.ListAllCharacters)
	r.GET("/characters/:id", characterController.GetCharacterByID)
}
