package repository

import (
	"github.com/jmoiron/sqlx"
)

type LikeRepository struct {
	db *sqlx.DB
}

func NewLikeRepository(db *sqlx.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

// SetLike добавляет или обновляет лайк/дизлайк (используем REPLACE или ручную проверку)
func (r *LikeRepository) SetLike(userID, targetID int, targetType string, value int) error {
	// Твоя структура таблицы имеет составной PRIMARY KEY (user_id, target_id, target_type)
	// Это позволяет использовать INSERT OR REPLACE
	query := `
		INSERT OR REPLACE INTO likes (user_id, target_id, target_type, value) 
		VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, userID, targetID, targetType, value)
	return err
}

// GetCounts возвращает количество лайков и дизлайков для объекта
func (r *LikeRepository) GetCounts(targetID int, targetType string) (likes int, dislikes int, err error) {
	query := `
		SELECT 
			COUNT(CASE WHEN value = 1 THEN 1 END),
			COUNT(CASE WHEN value = -1 THEN 1 END)
		FROM likes 
		WHERE target_id = ? AND target_type = ?`

	err = r.db.QueryRow(query, targetID, targetType).Scan(&likes, &dislikes)
	return
}

func (r *LikeRepository) DeleteLike(userID, targetID int, targetType string) error {
	query := `DELETE FROM likes WHERE user_id = ? AND target_id = ? AND target_type = ?`
	_, err := r.db.Exec(query, userID, targetID, targetType)
	return err
}

func (r *LikeRepository) GetCommentReactions(commentID int) (likes, dislikes int, err error) {
	query := `
		SELECT 
			COUNT(CASE WHEN value = 1 THEN 1 END),
			COUNT(CASE WHEN value = -1 THEN 1 END)
		FROM likes 
		WHERE target_id = ? AND target_type = 'comment'`
	err = r.db.QueryRow(query, commentID).Scan(&likes, &dislikes)
	return
}
