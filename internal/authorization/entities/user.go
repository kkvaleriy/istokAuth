package user

type user struct {
	Name     string
	Lastname string
	Nickname string
	Email    string
	Phone    int
	PassHash []byte
}
