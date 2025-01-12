package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetRepo struct {
	DB *sql.DB
}

func (m *SnippetRepo) Insert(
	title string, content string, expires int,
) (int, error) {
	stmnt := ` 
	    INSERT INTO snippets (title, content, created, expires)
        VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(),
        INTERVAL ? DAY)) 
    `

    result, err := m.DB.Exec(stmnt, title, content, expires)
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }


	return int(id), nil
}

func (m *SnippetRepo) Get(id int) (Snippet, error) {
	return Snippet{}, nil
}

func (m *SnippetRepo) Latest() ([]Snippet, error) {
	return nil, nil
}
