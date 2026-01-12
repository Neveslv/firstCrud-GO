package model

import "time"

type Usuario struct {
	ID        int       `json:"id"`
	Nome      string    `json:"nome"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Bio       string    `json:"bio"`
	Site      string    `json:"site"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}
type UsuarioInput struct {
	Nome     string `json:"nome" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
