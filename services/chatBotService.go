package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"porty-go/models"
	"porty-go/repositories"

	"go.mongodb.org/mongo-driver/mongo"
)

// Message represents a single message entry
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MessagesContainer holds an array of messages
type MessagesContainer struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type BotResponse struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int64  `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Usage             struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
		Message      struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GetServiceOpenAi() (models.IntegrationService, error) {
	serviceData, err := repositories.GetIntegrationServiceByName("OPENAI")
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.IntegrationService{}, errors.New("service not found")
		}
		return models.IntegrationService{}, err
	}

	return serviceData, nil
}

func GetServiceDialogFlow(idUser int, botService models.IntegrationService, newMessage string) (*string, error) {
	var logMessage []string // GET HISTORY DATA FROM DATABASE
	data := MessagesContainer{
		Messages: []Message{},
		Model:    botService.Model,
	}

	if len(logMessage) == 0 {
		for i := 1; i <= 2; i++ {
			var role = "system"
			if i%2 != 0 {
				role = "user"
			}
			data.Messages = append(data.Messages, Message{
				Role:    role,
				Content: newMessage,
			})
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", botService.ServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", botService.Token)) // Replace with actual token

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return nil, err
	}

	// Unmarshal response into struct
	var botResp BotResponse
	err = json.Unmarshal(body, &botResp)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return nil, err
	}

	content := botResp.Choices[0].Message.Content
	return &content, nil
}
