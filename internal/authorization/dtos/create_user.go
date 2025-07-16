package dtos

type CreateUserRequest struct {
	Name     string
	Lastname string
	Nickname string
	Email    string
	Phone    int
	Password string
}
