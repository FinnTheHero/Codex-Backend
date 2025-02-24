package chapter_service

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"Codex-Backend/api/internal/infrastructure/table"
	"errors"
	"slices"
)


func GetChapter(novel, chapter string) (domain.Chapter, error) {
	// novel = strings.ReplaceAll(novel, " ", "_")

	tables, err := table.GetTables()
	if err != nil {
		return domain.Chapter{}, errors.New("Failed to get tables")
	}

	if slices.Contains(tables, novel) {
		return repository.GetChapter(novel, chapter)
	}

	return domain.Chapter{}, errors.New("Novel not found")
}

func GetAllChapters(novel string) ([]domain.Chapter, error) {
	// novel = strings.ReplaceAll(novel, " ", "_")

	tables, err := table.GetTables()
	if err != nil {
		return nil, errors.New("Failed to get tables")
	}

	if slices.Contains(tables, novel) {
		return repository.GetAllChapters(novel)
	}

	return nil, errors.New("Novel not found")
}
