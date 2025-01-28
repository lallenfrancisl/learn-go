package mocks

import (
	"time"

	"github.com/lallenfrancisl/snippetbox/internal/models"
)

var mockSnippet = models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetRepo struct{}

func (m *SnippetRepo) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetRepo) Get(id int) (models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (m *SnippetRepo) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
