package usecase

import (
	"context"

	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
)

func (uc *userService) UpdateUserPassword(ctx context.Context, request *dtos.UpdateUserPasswordRequest) error {
	usr, err := user.UpdatePassword(request)
	if err != nil {
		return err
	}
	err = uc.repository.UpdateUserPassword(ctx, usr)
	if err != nil {
		return err
	}

	return nil
}
