package v1

import (
	"context"

	"github.com/kkvaleriy/istokAuthorization/internal/authorization/dtos"
)

type Usecase interface {
	SignUp(_ context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error)
}
