package service

import (
	cmn "Codex-Backend/api/common"
	firestore_client "Codex-Backend/api/internal/database/client"
	firestore_collections "Codex-Backend/api/internal/database/collections"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/timsims/pamphlet"
)

// Remove HTML tags and convert the body to Markdown
func cleanHtml(input string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(input))
	if err != nil {
		return "", err
	}

	doc.Find("head, script, style, nav, a, img, code, pre").Remove()

	html, err := doc.Find("body").Html()
	if err != nil {
		return "", err
	}

	return htmltomarkdown.ConvertString(html)
}

func CreateNovelFromEPUB(data []byte, ctx context.Context) error {
	parser, err := pamphlet.OpenBytes(data)
	if err != nil {
		return err
	}
	defer parser.Close()

	book := parser.GetBook()

	// Create Novel

	id, err := cmn.GenerateID("novel")
	if err != nil {
		return err
	}

	description, err := cleanHtml(book.Description)
	if err != nil {
		return err
	}

	novel := &domain.Novel{
		ID:          id,
		Title:       book.Title,
		Author:      book.Author,
		Description: description,
		CreatedAt:   cmn.TimeStamp(book.Date),
		UpdatedAt:   cmn.TimeStamp(""),
		Deleted:     false,
	}

	err, id = CreateNovel(*novel, ctx)
	if err != nil {
		return err
	}

	// Create chapters

	rawChapters := book.Chapters

	c_id, err := cmn.GenerateID("chapter")
	if err != nil {
		return err
	}

	chapters := make([]domain.Chapter, len(rawChapters))
	for i, chapter := range rawChapters {
		rawContent, err := chapter.GetContent()
		if err != nil {
			return err
		}

		content, err := cleanHtml(rawContent)
		if err != nil {
			return err
		}

		chapter := &domain.Chapter{
			ID:          c_id,
			Title:       chapter.Title,
			Author:      book.Author,
			Description: "",
			CreatedAt:   cmn.TimeStamp(""),
			UpdatedAt:   cmn.TimeStamp(""),
			Content:     content,
			Index:       i,
			Deleted:     false,
		}

		chapters[i] = *chapter
	}

	err = BatchUploadChapters(id, chapters, ctx)
	if err != nil {
		return err
	}

	return nil
}

func processChap(chapter pamphlet.Chapter, index int, author string) (*domain.Chapter, error) {
	c_id, err := cmn.GenerateID("chapter")
	if err != nil {
		return nil, err
	}

	rawContent, err := chapter.GetContent()
	if err != nil {
		return nil, err
	}

	titleLower := strings.ToLower(chapter.Title)
	contentLower := strings.ToLower(rawContent)

	if strings.HasPrefix(contentLower, titleLower) {
		rawContent = rawContent[len(chapter.Title):]
		rawContent = strings.TrimLeft(rawContent, " \t\n\r")
	}

	content, err := cleanHtml(rawContent)
	if err != nil {
		return nil, err
	}

	return &domain.Chapter{
		ID:          c_id,
		Title:       chapter.Title,
		Author:      author,
		Description: "",
		CreatedAt:   cmn.TimeStamp(""),
		UpdatedAt:   cmn.TimeStamp(""),
		Content:     content,
		Index:       index,
		Deleted:     false,
	}, nil
}

func CreateNovel(novel domain.Novel, ctx context.Context) (error, string) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err, ""
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	id, err := cmn.GenerateID("novel")
	if err != nil {
		return err, ""
	}

	novel.ID = id
	novel.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	novel.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	novel.Deleted = false

	err = c.CreateNovel(novel, ctx)
	if err != nil {
		return err, ""
	}

	return nil, id
}

func GetNovelById(id string, ctx context.Context) (*domain.Novel, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	novel, err := c.GetNovelById(id, ctx)
	if err != nil {
		return nil, err
	}

	if novel == nil {
		return nil, &cmn.Error{Err: errors.New("Novel Service Error - Get Novel - Novel with ID " + id + " not found"), Status: http.StatusNotFound}
	}

	return novel, nil
}

func GetNovelByTitle(title string, ctx context.Context) (*domain.Novel, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	novel, err := c.GetNovelByTitle(title, ctx)
	if err != nil {
		return nil, err
	}

	if novel == nil {
		return nil, &cmn.Error{Err: errors.New("Novel Service Error - Get Novel - Novel with title " + title + " not found"), Status: http.StatusNotFound}
	}

	return novel, nil
}

func GetAllNovels(ctx context.Context) (*[]domain.Novel, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	novels, err := c.GetAllNovels(ctx)
	if err != nil {
		return nil, err
	}

	if len(*novels) == 0 {
		return nil, &cmn.Error{Err: errors.New("Novel Service Error - Get All Novels - No novels found"), Status: http.StatusNotFound}
	}

	return novels, nil
}

func UpdateNovel(id string, novel domain.Novel, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	novel.ID = id
	novel.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	err = c.UpdateNovel(novel, ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteNovel(id string, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	err = c.DeleteNovel(id, ctx)
	if err != nil {
		return err
	}

	return nil
}
