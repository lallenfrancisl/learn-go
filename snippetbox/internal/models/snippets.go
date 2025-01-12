package models

import (
	"database/sql"
	"errors"
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
	stmnt := `
        SELECT id, title, content, created, expires FROM snippets
        WHERE expires > UTC_TIMESTAMP() AND id = ?
    `

	row := m.DB.QueryRow(stmnt, id)

	var s Snippet

	err := row.Scan(
		&s.ID, &s.Title,
		&s.Content, &s.Created,
		&s.Expires,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

func (m *SnippetRepo) Latest() ([]Snippet, error) {
	stmnt := `
        SELECT id, title, content, created, expires FROM snippets
        WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10
    `

	rows, err := m.DB.Query(stmnt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err := rows.Scan(
			&s.ID, &s.Title,
			&s.Content, &s.Created,
			&s.Expires,
		)

		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
