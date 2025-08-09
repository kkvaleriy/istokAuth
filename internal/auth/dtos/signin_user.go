package dtos

import (
	"time"

	"github.com/google/uuid"
)

// @Description User credentials.
type SignInRequest struct {
	Email    string `json:"email" example:"john@email.com" `
	Phone    int    `json:"phone" example:"79991112233"`
	Password string `json:"password" example:"mySuperPass" validate:"required,min=8"`
}

// @Description JWT.
type SignInResponse struct {
	JWT           string    `json:"jwt"`
	RToken        string    `json:"-"`
	ExpiresRToken time.Time `json:"-"`
}

type UUIDRequest struct {
	UUID uuid.UUID
}

func RequestByUUID(UUID string) (*UUIDRequest, error) {
	u, err := uuid.Parse(UUID)
	if err != nil {
		return nil, err
	}
	return &UUIDRequest{
		UUID: u,
	}, nil
}

type SignInError struct {
	Message string
	Reason  string
}

func (s *SignInError) Error() string {
	return s.Message
}
