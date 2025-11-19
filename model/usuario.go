package model

import "time"

type Usuario struct {
	ID        int       `json:"id"`
	Nome      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at"`
}
