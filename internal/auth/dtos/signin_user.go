package dtos

type SignINRequest struct {
	Email    string `json:"email"`
	Phone    int    `json:"phone"`
	Password string `json:"password"`
}
