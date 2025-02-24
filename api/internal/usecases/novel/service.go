package novel_service

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"errors"
)

type NovelService struct{}

func NewNovelService() *NovelService {
	return &NovelService{}
}

func (s *NovelService) GetNovel(novel string) (interface{}, error) {
	// title := strings.ReplaceAll(novel, " ", "_")

	NovelSchema, err := repository.GetNovel(novel)
	if err != nil {
		return nil, err
	}

	return NovelSchema, nil
}

func (s *NovelService) GetAllNovels() (interface{}, error) {
	result, err := repository.GetAllNovels()
	if err != nil {
		return nil, err
	}

	NovelDTOs, ok := result.([]domain.NovelDTO)
	if !ok {
		return nil, errors.New("Type assertion failed")
	}

	var novels []domain.Novel
	for _, dto := range NovelDTOs {
		novels = append(novels, dto.Novel)
	}

	return novels, nil
}
