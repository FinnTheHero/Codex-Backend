package services

import (
	aws_services "Codex-Backend/api/aws/services"
	"errors"
	"strings"
)

type ChapterService struct{}

func NewChapterService() *ChapterService {
	return &ChapterService{}
}

func (s *ChapterService) GetChapter(novel, chapter string) (interface{}, error) {
	novel = strings.ReplaceAll(novel, " ", "_")

	tables, err := aws_services.GetTables()
	if err != nil {
		return nil, errors.New("Failed to get tables")
	}

	for _, table := range tables {
		if table == novel {
			return aws_services.GetChapter(novel, chapter)
		}
	}

	return nil, errors.New("Novel not found")
}

func (s *ChapterService) GetAllChapters(novel string) (interface{}, error) {
	novel = strings.ReplaceAll(novel, " ", "_")

	tables, err := aws_services.GetTables()
	if err != nil {
		return nil, errors.New("Failed to get tables")
	}

	for _, table := range tables {
		if table == novel {
			return aws_services.GetAllChapters(novel)
		}
	}

	return nil, errors.New("Novel not found")
}
