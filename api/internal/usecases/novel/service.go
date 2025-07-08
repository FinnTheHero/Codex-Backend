package novel_service

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"Codex-Backend/api/internal/infrastructure/table"
	auth_service "Codex-Backend/api/internal/usecases/auth"
	"errors"
	"time"
)

func CreateNovel(novel domain.Novel) error {
	tableExists, err := table.IsTableCreated(novel.Title)
	if err != nil {
		return err
	}

	if tableExists {
		return errors.New("Novel already exists")
	}

	id, err := auth_service.GenerateID("novel")
	if err != nil {
		return err
	}

	novel.ID = id

	err = table.CreateTable(novel.ID)
	if err != nil {
		return err
	}

	novel.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	novel.UpdatedAt = novel.CreatedAt
	novel.UploadedAt = novel.CreatedAt

	err = repository.CreateNovel(novel)
	if err != nil {
		return err
	}

	return nil
}

func GetNovel(id string) (domain.Novel, error) {
	// title := strings.ReplaceAll(novel, " ", "_")

	novel, err := repository.GetNovel(id)
	if err != nil {
		return domain.Novel{}, err
	}

	return novel, nil
}

func GetAllNovels() (any, error) {
	result, err := repository.GetAllNovels()
	if err != nil {
		return nil, err
	}

	Novels, ok := result.([]domain.Novel)
	if !ok {
		return nil, errors.New("Type assertion failed")
	}

	var novels []domain.Novel
	for _, novel := range Novels {
		novels = append(novels, novel)
	}

	return novels, nil
}
