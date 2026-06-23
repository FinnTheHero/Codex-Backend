package service

import (
	cmn "Codex-Backend/api/common"
	db "Codex-Backend/api/internal/database"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"net/http"
	"strings"

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

	description, err := cleanHtml(book.Description)
	if err != nil {
		return err
	}

	// TODO: Fix in future commits: Add book creation time
	// createdAt, err := time.Parse(time.RFC3339, book.Date)
	// if err != nil {
	// 	return err
	// }

	newNovel := domain.Novel{
		Title:       book.Title,
		Author:      book.Author,
		Description: description,
	}

	// Create chapters

	rawChapters := book.Chapters

	// Split chapters by priority
	// Chpaters with chpater in title take priority
	// Notes, Synopsys and everything else will be processes last
	extraChapters := []pamphlet.Chapter{}
	actualChapters := []pamphlet.Chapter{}

	for _, chapter := range rawChapters {
		if strings.Contains(strings.ToLower(chapter.Title), "chapter") {
			actualChapters = append(actualChapters, chapter)
		} else {
			extraChapters = append(extraChapters, chapter)
		}
	}

	orderedChapters := append(actualChapters, extraChapters...)

	chapters := make([]domain.Chapter, len(orderedChapters))
	for i, chapter := range orderedChapters {
		chap, err := processChap(chapter, i, book.Author)
		if err != nil {
			return err
		}

		chapters[i] = *chap
	}

	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	err = client.CreateNovelFromEpub(newNovel, chapters, ctx)
	if err != nil {
		return err
	}

	return nil
}

func processChap(chapter pamphlet.Chapter, index int, author string) (*domain.Chapter, error) {
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
		Title:       chapter.Title,
		Author:      author,
		Description: "",
		Content:     content,
		Index:       index,
	}, nil
}

func CreateNovel(novel domain.CreateNovel, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if err = client.CreateNovel(novel, ctx); err != nil {
		return err
	}

	return nil
}

func GetNovelById(id string, ctx context.Context) (domain.Novel, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return domain.Novel{}, err
	}

	novel, err := client.GetNovelById(id, ctx)
	if err != nil {
		return domain.Novel{}, err
	}

	return novel, nil
}

func GetNovelByTitle(title string, ctx context.Context) (domain.Novel, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return domain.Novel{}, err
	}

	novel, err := client.GetNovelByTitle(title, ctx)
	if err != nil {
		return domain.Novel{}, err
	}

	return novel, nil
}

func GetAllNovels(ctx context.Context) ([]domain.Novel, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	novels, err := client.GetAllNovels(ctx)
	if err != nil {
		return nil, err
	}

	if len(novels) == 0 {
		return nil, &cmn.Error{Err: errors.New("Novel Service Error - Get All Novels - No novels found"), Status: http.StatusNotFound}
	}

	return novels, nil
}

func UpdateNovel(novel domain.Novel, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if err = client.UpdateNovel(novel, ctx); err != nil {
		return err
	}

	return nil
}

func DeleteNovel(id string, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if err = client.DeleteNovel(id, ctx); err != nil {
		return err
	}

	return nil
}
