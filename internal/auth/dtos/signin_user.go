package dtos

type SignInRequest struct {
	Email    string `json:"email"`
	Phone    int    `json:"phone"`
	Password string `json:"password"`
}
