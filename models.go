package main

import (
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/zrtgzrtg/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             database.Queries
	platform       string
}
type jsonError struct {
	Error string `json:"error"`
}
type requestBody struct {
	Body string `json:"body"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
}
type jsonClean struct {
	Cleaned_Body string `json:"cleaned_body"`
}
type requestChirpBody struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
