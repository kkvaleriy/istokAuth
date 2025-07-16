package usecase

import (
	"context"

	"github.com/kkvaleriy/istokAuthorization/internal/authorization/dtos"
	user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"
)

type UserCreator interface {
	AddUser(u *user.User) (*user.User, error)
}

type createUserUseCase struct {
	repository UserCreator
}

func NewCreateUser(repository UserCreator) *createUserUseCase {
	return &createUserUseCase{repository: repository}
}

func (uc *createUserUseCase) SignUP(_ context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	userForCreate, err := user.UserForSignUP(request)
	if err != nil {
		return nil, err
	}

	userForCreate, err = uc.repository.AddUser(userForCreate)
	if err != nil {
		return nil, err
	}

	response := &dtos.CreateUserResponse{
		UUID:      userForCreate.UUID,
		UserType:  userForCreate.UserType,
		CreatedAt: userForCreate.CreatedAt,
	}

	return response, nil
}
