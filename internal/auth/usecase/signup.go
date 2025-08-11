package usecase

import (
	"context"
	"fmt"

	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
)

func (uc *userService) SignUp(ctx context.Context, request *dtos.CreateUserRequest) (*dtos.CreateUserResponse, error) {
	userForCreate, err := user.SignUp(request)
	if err != nil {
		return nil, fmt.Errorf("signup: %w", err)
	}

	err = uc.repository.AddUser(ctx, userForCreate)
	if err != nil {
		return nil, fmt.Errorf("signup: %w", err)
	}

	response := &dtos.CreateUserResponse{
		UUID:      userForCreate.UUID,
		UserType:  userForCreate.UserType,
		CreatedAt: userForCreate.CreatedAt,
	}

	return response, nil
}
