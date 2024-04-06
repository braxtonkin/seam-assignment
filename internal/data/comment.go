package data

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

type Comment struct {
	// BlogID  int64  `json:"blogID"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

const ForeignKeyViolation = "23503"

var ForeignKeyError = errors.New("the blog post with the given id does not exist")

type CommentModel struct {
	DB *sql.DB
}

func (m CommentModel) Insert(comment *Comment, blogID int) (int, error) {
	var id int
	query := `
		INSERT INTO comments (blog_id, author, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := m.DB.QueryRow(query, blogID, comment.Author, comment.Content).Scan(&id)
	if err != nil {
		// tries to assert that error == pq.Error
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == ForeignKeyViolation {
			return -1, ForeignKeyError
		}
		return -1, err
	}
	return id, nil
}

func (m CommentModel) Get(id int) ([]Comment, error) {
	query := `
	SELECT author, content
	FROM comments
	WHERE blog_id = $1
	`

	var comments []Comment
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var input struct {
			author  string
			content string
		}

		if err := rows.Scan(&input.author, &input.content); err != nil {
			return nil, err
		}
		comments = append(comments, Comment{
			Author:  input.author,
			Content: input.content,
		})
	}

	return comments, nil
}
