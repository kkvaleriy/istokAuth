package dtos

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    int    `json:"phone"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	UUID uuid.UUID `json:"uuid"`
}
