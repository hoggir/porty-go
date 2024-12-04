package services

import (
	"porty-go/models"
	"porty-go/repositories"
)

type CharacterService struct {
	repo *repositories.CharacterRepository
}

func NewCharacterService(repo *repositories.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

func (s *CharacterService) GetAllCharacters() ([]models.Character, error) {
	return s.repo.GetAllCharacters()
}
