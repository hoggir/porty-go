package controllers

import (
	"net/http"
	"porty-go/models"
	"porty-go/services"

	"github.com/gin-gonic/gin"
)

type CharacterController struct {
	service *services.CharacterService
}

func NewCharacterController(service *services.CharacterService) *CharacterController {
	return &CharacterController{service: service}
}

// GetAllCharacters godoc
// @Summary Get all characters
// @Description Get all characters from the database
// @Tags characters
// @Produce json
// @Success 200 {array} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters [get]
func (cc *CharacterController) GetAllCharacters(c *gin.Context) {
	characters, err := cc.service.GetAllCharacters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    characters,
	})
}
