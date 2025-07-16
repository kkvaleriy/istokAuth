package v1

import (
	"net/http"

	"github.com/kkvaleriy/istokAuthorization/internal/authorization/dtos"
	"github.com/labstack/echo/v4"
)

type handler struct {
	usecase Usecase
}

func NewHandler(uc Usecase) *handler {
	return &handler{usecase: uc}
}

func (h *handler) signUp(c echo.Context) error {

	var request *dtos.CreateUserRequest

	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "bad json in request")
	}

	response, err := h.usecase.SignUp(c.Request().Context(), request)
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, response)
}
