package dtos

type SignInRequest struct {
	Email    string `json:"email" example:"john@email.com" validate:"email"`
	Phone    int    `json:"phone example:"79991112233""`
	Password string `json:"password example:"mySuperPass" validate:"required,min=8"`
}
