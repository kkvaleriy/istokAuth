package usecase

import (
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
	"golang.org/x/net/context"
)

type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Repository interface {
	AddUser(ctx context.Context, u *user.User) error
}

type usecase struct {
	repository Repository
	log        logger
}

func New(repository Repository, log logger) *usecase {
	return &usecase{repository: repository,
		log: log}
}
