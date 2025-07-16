package usecase

import (
	user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"
	"golang.org/x/net/context"
)

type Repository interface {
	AddUser(ctx context.Context, u *user.User) error
}

type createUserUseCase struct {
	repository Repository
}

func NewCreateUser(repository Repository) *createUserUseCase {
	return &createUserUseCase{repository: repository}
}
