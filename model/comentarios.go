package model

import "time"

type Comentarios struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id" binding:"required"`
	UserID    int       `json:"user_id" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}
