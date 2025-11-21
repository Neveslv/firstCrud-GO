package model

import "time"

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id" binding:"required"`
	Titulo    string    `json:"title" binding:"required"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	// Campos extras para quando retornar com dados do usuário
	UserName  string `json:"user_name,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
}

type PostInput struct {
	Titulo  string `json:"titulo"binding:"required,min=3,max=200"`
	Content string `json:"content"binding:"required,min=10"`
}
