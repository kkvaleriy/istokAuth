package v1

import (
	_ "github.com/kkvaleriy/istokAuth/docs"
	"github.com/kkvaleriy/istokAuth/internal/auth/delivery/http/v1/middleware"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *handler) Routes(domain *echo.Group, secret string) {
	// Swagger endpoint
	domain.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := domain.Group("/auth")
	auth.POST("/signup", h.signUp)
	auth.POST("/signin", h.signIn)
	auth.GET("/refresh", h.Refresh)

	protected := domain.Group("/user", middleware.JWTAuthCheck([]byte(secret)))
	protected.PUT("/update-password", h.UpdateUserPassword)
}
