package model

import "time"

type Comentarios struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id" binding:"required"`
	UserID    int       `json:"user_id" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at"`

	UserNome  string `json:"user_nome,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

type ComentarioInput struct {
	Content string `json:"content" binding:"required,min=1 max=1000"`
}
