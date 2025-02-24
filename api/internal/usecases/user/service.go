package user_service

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"Codex-Backend/api/internal/infrastructure/table"
	"errors"
	"strings"
)


func CreateNovel(novel domain.Novel) error {
	tableExists, err := table.IsTableCreated(novel.Title)
	if err != nil {
		return err
	}

	if tableExists {
		return errors.New("Novel already exists")
	}

	err = table.CreateTable(novel.Title)
	if err != nil {
		return err
	}

	err = repository.CreateNovel(novel)
	if err != nil {
		return err
	}

	return nil
}

func CreateChapter(novel string, chapter domain.Chapter) error {
	novel = strings.ReplaceAll(novel, " ", "_")

	tableExists, err := table.IsTableCreated(novel)
	if err != nil {
		return err
	}

	if !tableExists {
		return errors.New("Novel not found")
	}

	err = repository.CreateChapter(novel, chapter)
	if err != nil {
		return err
	}

	return nil
}
