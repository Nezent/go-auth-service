package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}
