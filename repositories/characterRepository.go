package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"porty-go/models"

	"github.com/supabase-community/supabase-go"
)

type CharacterRepository struct {
	client *supabase.Client
}

func NewCharacterRepository() (*CharacterRepository, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	client, err := supabase.NewClient(supabaseURL, supabaseKey, &supabase.ClientOptions{})

	if err != nil {
		return nil, err
	}
	return &CharacterRepository{client: client}, nil
}

func (r *CharacterRepository) GetAllCharacters(page, record int, search string) ([]models.Character, error) {
	var characters []models.Character

	// Define the RPC parameters
	params := map[string]interface{}{
		"search_query": search,
		"page":         page,
		"record":       record,
	}

	// Call the PostgreSQL function via RPC
	resp := r.client.Rpc("list_characters", "", params)

	// Parse the response
	err := json.Unmarshal([]byte(resp), &characters)

	// Check for errors
	if err != nil {
		var errorResponse models.SupRpcErrorResponse
		err := json.Unmarshal([]byte(resp), &errorResponse)
		// Set default error message if parsing fails
		if err != nil {
			return nil, err
		}
		return nil, errors.New(errorResponse.Message)
	}

	return characters, nil
}
