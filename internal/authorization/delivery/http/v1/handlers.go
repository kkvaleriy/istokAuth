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

type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type handler struct {
	usecase Usecase
	log     logger
}

func NewHandler(uc Usecase, log logger) *handler {
	return &handler{usecase: uc, log: log}
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
