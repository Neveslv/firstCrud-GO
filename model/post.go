package model

import "time"

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" binding:"required"`
	Titulo    string    `json:"title" binding:"required"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	// Campos extras para JOIN
	UserName string `json:"user_name,omitempty"`
}
