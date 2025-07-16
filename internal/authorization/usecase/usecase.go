package usecase

import user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"

type Repository interface {
	AddUser(u *user.User) (*user.User, error)
}

type createUserUseCase struct {
	repository Repository
}

func NewCreateUser(repository Repository) *createUserUseCase {
	return &createUserUseCase{repository: repository}
}
