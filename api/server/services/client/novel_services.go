package client_services

import (
	aws_services "Codex-Backend/api/aws/services"
	"Codex-Backend/api/models"
	"errors"
)

type NovelService struct{}

func NewNovelService() *NovelService {
	return &NovelService{}
}

func (s *NovelService) GetNovel(novel string) (interface{}, error) {
	// title := strings.ReplaceAll(novel, " ", "_")

	NovelSchema, err := aws_services.GetNovel(novel)
	if err != nil {
		return nil, err
	}

	return NovelSchema, nil
}

func (s *NovelService) GetAllNovels() (interface{}, error) {
	result, err := aws_services.GetAllNovels()
	if err != nil {
		return nil, err
	}

	NovelDTOs, ok := result.([]models.NovelDTO)
	if !ok {
		return nil, errors.New("Type assertion failed")
	}

	var novels []models.Novel
	for _, dto := range NovelDTOs {
		novels = append(novels, dto.Novel)
	}

	return novels, nil
}
