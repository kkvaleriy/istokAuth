package usecase

import (
	"context"

	"github.com/kkvaleriy/istokAuthorization/internal/authorization/dtos"
	user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"
)

func (uc *createUserUseCase) SignUp(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	userForCreate, err := user.UserForSignUp(request)
	if err != nil {
		return nil, err
	}

	err = uc.repository.AddUser(ctx, userForCreate)
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
