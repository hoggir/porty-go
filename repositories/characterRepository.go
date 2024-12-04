package repositories

import (
	"encoding/json"
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

func (r *CharacterRepository) GetAllCharacters() ([]models.Character, error) {
	resp, _, err := r.client.From("characters").Select("*", "", false).Execute()
	if err != nil {
		return nil, err
	}

	var characters []models.Character
	if err := json.Unmarshal(resp, &characters); err != nil {
		return nil, err
	}

	return characters, nil
}
