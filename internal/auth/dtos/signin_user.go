package dtos

import "time"

// @Description User credentials.
type SignInRequest struct {
	Email    string `json:"email" example:"john@email.com" validate:"email"`
	Phone    int    `json:"phone" example:"79991112233"`
	Password string `json:"password" example:"mySuperPass" validate:"required,min=8"`
}

// @Description JWT.
type SignInResponse struct {
	JWT           string `json:"jwt"`
	RToken        string
	ExpiresRToken time.Time
}

type SignInError struct {
	Message string
	Reason  string
}

func (s *SignInError) Error() string {
	return s.Message
}
