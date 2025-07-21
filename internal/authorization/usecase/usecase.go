package usecase

import (
	user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"
	"golang.org/x/net/context"
)

type Repository interface {
	AddUser(ctx context.Context, u *user.User) error
}

type usecase struct {
	repository Repository
}

func New(repository Repository) *usecase {
	return &usecase{repository: repository}
}
