package usecase

import (
	"context"
	"fmt"

	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
)

func (uc *userService) UpdateUserPassword(ctx context.Context, request *dtos.UpdateUserPasswordRequest) error {
	usr, err := user.UpdatePassword(request)
	if err != nil {
		return fmt.Errorf("update user password: %w", err)
	}
	err = uc.repository.UpdateUserPassword(ctx, usr)
	if err != nil {
		return fmt.Errorf("update user password: %w", err)
	}

	return nil
}
