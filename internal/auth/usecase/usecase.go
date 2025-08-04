package usecase

import (
	"time"

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
	CheckUserByCredentials(ctx context.Context, u *user.User) (*user.User, error)
	AddToken(ctx context.Context, t *user.RToken) error
}

type userService struct {
	secret     string
	jwtTTL     time.Duration
	refreshTTL time.Duration
	repository Repository
	log        logger
}

func NewUserService(secret string, repository Repository, log logger) *userService {
	return &userService{repository: repository,
		log: log}
}
