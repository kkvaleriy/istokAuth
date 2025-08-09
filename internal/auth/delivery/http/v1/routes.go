package v1

import (
	_ "github.com/kkvaleriy/istokAuth/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (h *handler) Routes(domain *echo.Group) {
	// Swagger endpoint
	domain.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := domain.Group("/auth")
	auth.POST("/signup", h.signUp)
	auth.POST("/signin", h.signIn)
}
