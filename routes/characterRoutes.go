package routes

import (
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
		panic(err)
	}
	characterController := controllers.NewCharacterController(services.NewCharacterService(repo))

	r.GET("/characters", characterController.GetAllCharacters)
}
