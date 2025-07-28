package v1

import (
	"context"

	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
)

type Usecase interface {
	SignUp(_ context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error)
}
