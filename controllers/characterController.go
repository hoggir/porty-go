package controllers

import (
	"net/http"
	"porty-go/models"
	"porty-go/services"
	"strconv"

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
// @Param page query int false "Page number"
// @Param record query int false "Number of records per page"
// @Param search query string false "Search term"
// @Success 200 {array} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters [get]
func (cc *CharacterController) ListAllCharacters(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	record, _ := strconv.Atoi(c.DefaultQuery("record", "10"))
	search := c.DefaultQuery("search", "")

	characters, err := cc.service.ListAllCharacters(page, record, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    characters,
	})
}

// GetCharacterByID godoc
// @Summary Get character by ID
// @Description Get a character by ID from the database
// @Tags characters
// @Produce json
// @Param id path int true "Character ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /characters/{id} [get]
func (cc *CharacterController) GetCharacterByID(c *gin.Context) {
	id := c.Param("id")

	character, err := cc.service.GetCharacterByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Character retrieved successfully",
		Data:    character,
	})
}
