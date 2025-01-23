package services

import (
	"fmt"
	"os"
	"porty-go/models"
	"porty-go/repositories"
	"porty-go/utils"
)

type CharacterService struct {
	repo *repositories.CharacterRepository
}

func NewCharacterService(repo *repositories.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

func (s *CharacterService) ListAllCharacters(page, record int, search string) ([]models.Character, error) {
	return s.repo.GetAllCharacters(page, record, search)
}

func (s *CharacterService) GetCharacterByID(id string) (models.Character, error) {
	character, err := s.repo.GetCharacterByID(id)
	if err != nil {
		return models.Character{}, err
	}
	encryptKey := os.Getenv("ENCRYPT_KEY")
	encryptedID, err := utils.Encrypt(fmt.Sprintf("%d", *character.ID), encryptKey)
	if err != nil {
		return models.Character{}, err
	}
	character.ID = nil
	character.EncryptedID = &encryptedID
	return character, nil
}
