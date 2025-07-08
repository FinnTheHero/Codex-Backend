package chapter_service

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"Codex-Backend/api/internal/infrastructure/table"
	auth_service "Codex-Backend/api/internal/usecases/auth"
	"errors"
	"slices"
	"strings"
)

func CreateChapter(novelId string, chapter domain.Chapter) error {
	novelId = strings.ReplaceAll(novelId, " ", "_")

	tableExists, err := table.IsTableCreated(novelId)
	if err != nil {
		return err
	}

	if !tableExists {
		return errors.New("Novel not found")
	}

	id, err := auth_service.GenerateID("chapter")
	if err != nil {
		return err
	}

	chapter.ID = id

	err = repository.CreateChapter(novelId, chapter)
	if err != nil {
		return err
	}

	return nil
}

func GetChapter(novelId, chapterId string) (domain.Chapter, error) {
	// novel = strings.ReplaceAll(novel, " ", "_")

	tableIds, err := table.GetTableIds()
	if err != nil {
		return domain.Chapter{}, errors.New("Failed to get tables")
	}

	if slices.Contains(tableIds, novelId) {
		return repository.GetChapter(novelId, chapterId)
	}

	return domain.Chapter{}, errors.New("Novel not found")
}

func GetAllChapters(novelId string) ([]domain.Chapter, error) {
	// novel = strings.ReplaceAll(novel, " ", "_")

	tableIds, err := table.GetTableIds()
	if err != nil {
		return nil, errors.New("Failed to get tables")
	}

	if slices.Contains(tableIds, novelId) {
		return repository.GetAllChapters(novelId)
	}

	return nil, errors.New("Novel not found")
}
