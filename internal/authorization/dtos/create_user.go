package dtos

import (
	"time"

	"github.com/google/uuid"
)

// @Description User account information for registration
type CreateUserRequest struct {
	Name     string `json:"name" example:"John"`
	Lastname string `json:"lastname" example:"Doe"`
	Nickname string `json:"nickname" example:"Johny1"`
	Email    string `json:"email" example:"john@email.com"`
	Phone    int    `json:"phone" example:"79990001122"`
	Password string `json:"password" example:"mySuperPass"`
}

// @Description Information about the user's account after successful registration
type CreateUserResponse struct {
	UUID      uuid.UUID `json:"uuid" example:"16763be4-6022-406e-a950-fcd5018633ca"`
	UserType  string    `json:"user_type" example:"USER"`
	CreatedAt time.Time `json:"created_at" example:"2006-01-02T03:04:05.5141511+03:00"`
}

type ValidationError struct {
	Message string
	Field   string
	Value   string
}

func (v ValidationError) Error() string {
	return v.Message
}
