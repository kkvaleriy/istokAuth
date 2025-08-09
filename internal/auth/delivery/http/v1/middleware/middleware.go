package middleware

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	httperrors "github.com/kkvaleriy/istokAuth/internal/auth/delivery/http/v1/errors"
	"github.com/labstack/echo/v4"
)

func JWTAuthCheck(secret []byte) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return &httperrors.AuthError{Err: errors.New("missing Authorization header")}
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return &httperrors.AuthError{Err: errors.New("invalid Authorization format")}
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return secret, nil
			})
			if err != nil || !token.Valid {
				return &httperrors.AuthError{Err: errors.New("invalid or expired token")}
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return &httperrors.AuthError{Err: errors.New("invalid token claims")}
			}

			c.Set("user", claims)
			return next(c)
		}
	}
}
