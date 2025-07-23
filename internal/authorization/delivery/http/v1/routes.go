package v1

import "github.com/labstack/echo/v4"

func (h *handler) Routes(domain *echo.Group) {

	domain.POST("/signup", h.signUp)
	h.log.Debug("added route for /signup")
}
