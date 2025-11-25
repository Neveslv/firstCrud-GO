package model

import "time"

type Usuario struct {
	ID        int       `json:"id"`
	Nome      string    `json:"nome"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
type UsuarioInput struct {
	Nome  string `json:"nome" binding:"required,min=3,max=100"`
	Email string `json:"email" binding:"required,email"`
}
