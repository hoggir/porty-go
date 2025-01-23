package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (r *CharacterRepository) GetCharacterByID(id string) (models.Character, error) {
	var characters models.Character

	// Fetch a single record from the "characters"
	resp, _, err := r.client.From("characters").Select("*", "exact", false).Eq("id", id).Single().Execute()
	if err != nil {
		fmt.Println("Error executing query: ", err)
		if err.Error() == "(PGRST116) JSON object requested, multiple (or no) rows returned" {
			return characters, errors.New("character not found")
		}
		return characters, err
	}

	// Parse the JSON response into the characters slice
	err = json.Unmarshal(resp, &characters)
	if err != nil {
		fmt.Println("Error parsing response: ", err)
		return characters, err
	}

	return characters, nil
}
