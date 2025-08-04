package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
)

func (uc *userService) SignIn(ctx context.Context, request *dtos.SignInRequest) (*dtos.SignInResponse, error) {
	u, err := user.SignIn(request)
	if err != nil {
		return nil, err
	}

	u, err = uc.repository.CheckUserByCredentials(ctx, u)
	if err != nil {
		return nil, err
	}

	if !u.IsActive {
		return nil, fmt.Errorf("Invalid credentials")
	}

	jwtToken, err := uc.GenerateJWT(u)
	if err != nil {
		return nil, fmt.Errorf("Internal error")
	}

	rToken := u.RefreshToken(uc.refreshTTL)

	err = uc.repository.AddToken(ctx, rToken)
	if err != nil {
		return nil, err
	}

	return &dtos.SignInResponse{JWT: jwtToken, RToken: rToken.UUID.String()}, nil
}

func (uc *userService) GenerateJWT(u *user.User) (string, error) {
	issuedAt := time.Now().Unix()
	expiresAt := time.Unix(issuedAt, 0).Add(uc.jwtTTL).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  "IstokAuth",
		"sub":  u.UUID,
		"role": u.UserType,
		"exp":  expiresAt,
		"iat":  issuedAt,
		"nbf":  issuedAt,
	})

	return token.SignedString(uc.secret)
}
