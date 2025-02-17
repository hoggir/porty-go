package controllers

import (
	"net/http"
	"porty-go/models"
	"porty-go/services"

	"github.com/gin-gonic/gin"
)

type TestAiBody struct {
	Message string `json:"message"`
}

// TestAI godoc
// @Summary Test Chat Bot AI
// @Description Just Chatbot AI testing
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body TestAiBody true "Message"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /chat [post]
func ChatAi(c *gin.Context) {
	var messageBody TestAiBody
	if err := c.ShouldBindJSON(&messageBody); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	if messageBody.Message == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Status:  "error",
			Message: "Message is required!",
		})
		return
	}

	aiService, err := services.GetServiceOpenAi()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	chatResponse, err := services.GetServiceDialogFlow(1, aiService, messageBody.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status:  "success",
		Message: "Message received",
		Data:    chatResponse,
	})
}
