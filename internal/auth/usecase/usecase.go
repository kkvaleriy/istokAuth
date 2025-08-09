package usecase

import (
	"time"

	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
	"golang.org/x/net/context"
)

type Authentificator interface {
	SignUp(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error)
	SignIn(ctx context.Context, request *dtos.SignInRequest) (*dtos.SignInResponse, error)
	UpdateUserPassword(ctx context.Context, request *dtos.UpdateUserPasswordRequest) error
	Refresh(ctx context.Context, request *dtos.UUIDRequest) (*dtos.SignInResponse, error)
}

type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Repository interface {
	AddUser(ctx context.Context, u *user.User) error
	UpdateUserPassword(ctx context.Context, u *user.User) error
	CheckUserByCredentials(ctx context.Context, u *user.User) (*user.User, error)
	AddToken(ctx context.Context, t *user.RToken) error
	RefreshToken(ctx context.Context, u *user.User, t *user.RToken) (*user.User, error)
}

type tokenConfigurator interface {
	SecretKey() string
	RefreshTTL() time.Duration
	AccessTTL() time.Duration
}

type tokenConfig struct {
	Secret     []byte
	RefreshTTL time.Duration
	AccessTTL  time.Duration
}

type userService struct {
	repository Repository
	token      *tokenConfig
	log        logger
}

func NewUserService(tParams tokenConfigurator, repository Repository, log logger) *userService {
	return &userService{repository: repository,
		token: &tokenConfig{
			Secret:     []byte(tParams.SecretKey()),
			RefreshTTL: tParams.RefreshTTL(),
			AccessTTL:  tParams.AccessTTL(),
		},
		log: log}
}
