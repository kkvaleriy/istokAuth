package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
)

func (uc *userService) Refresh(ctx context.Context, request *dtos.UUIDRequest) (*dtos.SignInResponse, error) {

	rToken := user.TokenRefresh(request.UUID)
	usr := user.Empty()

	usr, err := uc.repository.RefreshToken(ctx, usr, rToken)
	if err != nil {
		return nil, fmt.Errorf("refresh: %w", err)
	}

	rToken = usr.RefreshToken(uc.token.RefreshTTL)

	err = uc.repository.AddToken(ctx, rToken)
	if err != nil {
		return nil, fmt.Errorf("refresh: %w", err)
	}

	aToken, err := uc.GenerateJWT(usr)
	if err != nil {
		uc.log.Error("Can'n generate JWT for user", "user", usr.Nickname)
		return nil, fmt.Errorf("refresh: %w", err)
	}

	return &dtos.SignInResponse{JWT: aToken, RToken: rToken.UUID.String(),
		ExpiresRToken: time.Unix(rToken.ExpiresAt, 0)}, nil
}
