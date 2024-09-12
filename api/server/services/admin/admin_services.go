package admin_services

import (
	aws_services "Codex-Backend/api/aws/services"
	"Codex-Backend/api/models"
	"errors"
	"strings"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) CreateNovel(novel models.Novel) error {
	err := aws_services.CreateTable(novel.Title)
	if err != nil {
		return err
	}

	err = aws_services.CreateNovel(novel)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminService) CreateChapter(novel string, chapter models.Chapter) error {
	novel = strings.ReplaceAll(novel, " ", "_")

	tableExists, err := aws_services.IsTableCreated(novel)
	if err != nil {
		return err
	}

	if !tableExists {
		return errors.New("Novel not found")
	}

	err = aws_services.CreateChapter(novel, chapter)
	if err != nil {
		return err
	}

	return nil
}
