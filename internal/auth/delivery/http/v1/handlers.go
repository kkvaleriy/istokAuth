package v1

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	"github.com/kkvaleriy/istokAuth/internal/auth/usecase"
	"github.com/labstack/echo/v4"
)

var (
	ErrNotUniqUser = errors.New("not uniq user.")
)

type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type handler struct {
	usecase usecase.Authentificator
	log     logger
}

func NewHandler(uc usecase.Authentificator, log logger) *handler {
	return &handler{usecase: uc, log: log}
}

// @Summary Registering a new user
// @Tags Authorization
// @Description Registering a new user
// @Accept json
// @Produce json
// @Param input body dtos.CreateUserRequest true "Account information for signup"
// @Success 200 {object} dtos.CreateUserResponse "Information about the user's account after successful registration"
// @Failure 400 {object} badRequestErrorResponse "Bad request"
// @Failure 409 {object} validationDTOErrorResponse "A user already exists"
// @Failure 422 {object} validationErrorResponse "Bad json in request"
// @Failure 500 {object} internalServerErrorResponse "Internal server error"
// @Router /signup [post]
func (h *handler) signUp(c echo.Context) error {

	request := &dtos.CreateUserRequest{}

	if err := c.Bind(request); err != nil {
		body := []byte{}

		c.Request().Body.Read(body)
		defer c.Request().Body.Close()
		h.log.Error("can't parse json in request", "host", c.Request().Host, "URL", c.Request().URL, "body", body, "error", err)

		return &BadRequestError{Err: err}
	}

	h.log.Debug("request for signup", "host", c.Request().Host, "URL", c.Request().URL, "request", request)

	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return ErrValidation(err)
	}

	response, err := h.usecase.SignUp(c.Request().Context(), request)
	if err != nil {
		var validationError *dtos.ValidationError
		if errors.As(err, &validationError) {
			return &ValidationDTOError{Err: validationError}
		}

		return err
	}

	h.log.Info("new user has been created", "host", c.Request().Host, "URL", c.Request().URL, "request", request, "result", response)

	return c.JSON(http.StatusCreated, response)
}
