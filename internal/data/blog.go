package data

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

// Consider adding a Created At field (of type Time)
type Blog struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

type BlogModel struct {
	DB *sql.DB
}

var RecordNotFound = errors.New("record not found")

func (m BlogModel) Insert(blog *Blog) (int, error) {
	var id int
	query := `
		INSERT INTO blogs (title, author, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := m.DB.QueryRow(query, blog.Title, blog.Author, blog.Content).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// Retrieves a blog post associated with the given ID from the database
func (m BlogModel) Get(id int) (Blog, error) {
	if id < 1 {
		return Blog{}, RecordNotFound
	}

	var blog Blog

	query := `
	SELECT title, author, content
	FROM blogs
	WHERE id = $1
	`

	err := m.DB.QueryRow(query, id).Scan(
		&blog.Title,
		&blog.Author,
		&blog.Content,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Blog{}, RecordNotFound
		default:
			return Blog{}, err
		}
	}

	return blog, nil
}

func (m BlogModel) GetAll() ([]Blog, error) {
	query := `
	SELECT title, author, content
	FROM blogs
	`

	var blogs []Blog
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var input struct {
			title   string
			author  string
			content string
		}

		if err := rows.Scan(&input.title, &input.author, &input.content); err != nil {
			return nil, err
		}
		blogs = append(blogs, Blog{
			Title:   input.title,
			Author:  input.author,
			Content: input.content,
		})
	}

	return blogs, nil

}

func (m BlogModel) Update(blog *Blog, id int) error {
	query := ` UPDATE blogs
	SET title = $1, author = $2, content = $3
	WHERE id = $4
	`

	_, err := m.DB.Exec(query, blog.Title, blog.Content, blog.Author, id)
	if err != nil {
		return err
	}
	return nil
}

func (m BlogModel) Delete(id int) error {
	query := ` DELETE FROM blogs
	WHERE id = $1
	`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return RecordNotFound
	}

	return nil
}
