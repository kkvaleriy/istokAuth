package usecase

import (
	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
	"golang.org/x/net/context"
)

type Authentificator interface {
	SignUp(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error)
}

type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Repository interface {
	AddUser(ctx context.Context, u *user.User) error
}

type userService struct {
	repository Repository
	log        logger
}

func New(repository Repository, log logger) *userService {
	return &userService{repository: repository,
		log: log}
}
