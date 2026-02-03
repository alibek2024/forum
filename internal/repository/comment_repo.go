package repository

import (
	"github.com/alibek2024/forum/internal/models"
	"github.com/jmoiron/sqlx"
)

type CommentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(comment *models.Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, comment.PostID, comment.UserID, comment.Content)
	return err
}

func (r *CommentRepository) GetByPostID(postID int) ([]models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at 
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC`

	rows, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.AuthorName, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

// UPDATE: Редактировать комментарий
func (r *CommentRepository) Update(commentID, userID int, newContent string) error {
	query := `UPDATE comments SET content = ? WHERE id = ? AND user_id = ?`
	_, err := r.db.Exec(query, newContent, commentID, userID)
	return err
}

// DELETE: Удалить комментарий
func (r *CommentRepository) Delete(commentID, userID int) error {
	query := `DELETE FROM comments WHERE id = ? AND user_id = ?`
	_, err := r.db.Exec(query, commentID, userID)
	return err
}
func (r *CommentRepository) GetByID(id int) (*models.Comment, error) {
	query := `SELECT id, post_id, user_id, content, created_at FROM comments WHERE id = ?`
	c := &models.Comment{}
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}
