package client_services

import (
	aws_services "Codex-Backend/api/aws/services"
	"strings"
)

type NovelService struct{}

func NewNovelService() *NovelService {
	return &NovelService{}
}

func (s *NovelService) GetNovel(novel string) (interface{}, error) {
	title := strings.ReplaceAll(novel, " ", "_")

	NovelSchema, err := aws_services.GetNovel(title)
	if err != nil {
		return nil, err
	}

	return NovelSchema, nil
}

func (s *NovelService) GetAllNovels() (interface{}, error) {
	Novels, err := aws_services.GetAllNovels()
	if err != nil {
		return nil, err
	}

	return Novels, nil
}
