package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
)

func (uc *userService) SignIn(ctx context.Context, request *dtos.SignInRequest) (*dtos.SignInResponse, error) {
	usr, err := user.SignIn(request)
	if err != nil {
		return nil, fmt.Errorf("signin: %w")
	}

	usr, err = uc.repository.CheckUserByCredentials(ctx, usr)
	if err != nil {
		return nil, fmt.Errorf("signin: %w")
	}

	if !usr.IsActive {
		return nil, errors.New("Invalid credentials")
	}

	jwtToken, err := uc.GenerateJWT(usr)
	if err != nil {
		return nil, errors.New("Internal error")
	}

	rToken := usr.RefreshToken(uc.token.RefreshTTL)

	err = uc.repository.AddToken(ctx, rToken)
	if err != nil {
		return nil, fmt.Errorf("signin: %w")
	}

	return &dtos.SignInResponse{JWT: jwtToken, RToken: rToken.UUID.String(),
		ExpiresRToken: time.Unix(rToken.ExpiresAt, 0)}, nil
}

func (uc *userService) GenerateJWT(u *user.User) (string, error) {
	issuedAt := time.Now().Unix()
	expiresAt := time.Unix(issuedAt, 0).Add(uc.token.AccessTTL).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "IstokAuth",
		"sub":  u.UUID,
		"role": u.UserType,
		"exp":  expiresAt,
		"iat":  issuedAt,
		"nbf":  issuedAt,
	})

	return token.SignedString(uc.token.Secret)
}
