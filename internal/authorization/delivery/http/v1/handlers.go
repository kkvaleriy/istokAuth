package v1

import (
	"errors"
	"net/http"

	"github.com/kkvaleriy/istokAuthorization/internal/authorization/dtos"
	"github.com/labstack/echo/v4"
)

var (
	ErrNotUniqUser = errors.New("Not uniq user:")
)

type handler struct {
	usecase Usecase
}

func NewHandler(uc Usecase) *handler {
	return &handler{usecase: uc}
}

func (h *handler) signUp(c echo.Context) error {

	request := &dtos.CreateUserRequest{}

	if err := c.Bind(request); err != nil {
		return echo.NewHTTPError(echo.ErrBadRequest.Code, "bad json in request")
	}

	response, err := h.usecase.SignUp(c.Request().Context(), request)
	if err != nil {
		if errors.Is(err, ErrNotUniqUser) {
			return echo.ErrConflict
		}
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, response)
}
