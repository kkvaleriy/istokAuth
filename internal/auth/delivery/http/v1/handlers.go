package v1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
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
	usecase Usecase
	log     logger
}

func NewHandler(uc Usecase, log logger) *handler {
	return &handler{usecase: uc, log: log}
}

// @Summary Registering a new user
// @Tags Authorization
// @Description Registering a new user
// @Accept json
// @Produce json,plain
// @Param input body dtos.CreateUserRequest true "Account information for signup"
// @Success 200 {object} dtos.CreateUserResponse "Information about the user's account after successful registration"
// @Failure 409 {string} string "A user already exists"
// @Failure 400 {string} string "Bad json in request"
// @Failure 500 {string} string "Internal server error"
// @Router /signup [post]
func (h *handler) signUp(c echo.Context) error {

	request := &dtos.CreateUserRequest{}

	if err := c.Bind(request); err != nil {
		body := []byte{}

		c.Request().Body.Read(body)
		defer c.Request().Body.Close()
		h.log.Error("can't parse json in request", "host", c.Request().Host, "URL", c.Request().URL, "body", body, "error", err)

		return c.String(http.StatusBadRequest, "bad json in request")
	}

	h.log.Debug("request for signup", "host", c.Request().Host, "URL", c.Request().URL, "request", request)

	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		failedValidate := "fields that failed to validate: "
		for _, err := range err.(validator.ValidationErrors) {
			failedValidate = fmt.Sprintf("%s%s:%s:%s ", failedValidate, err.Field(), err.Tag(), err.Error())
		}
		return c.String(http.StatusBadRequest, failedValidate)
	}

	response, err := h.usecase.SignUp(c.Request().Context(), request)
	if err != nil {
		var validationError *dtos.ValidationError
		if errors.As(err, &validationError) {
			h.log.Error("not uniq user", "host", c.Request().Host, "URL", c.Request().URL, "request", request, "error", validationError.Error())

			return c.String(http.StatusConflict, err.Error())
		}

		h.log.Error("unexpected error", "error", err)

		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	h.log.Info("new user has been created", "host", c.Request().Host, "URL", c.Request().URL, "request", request, "result", response)

	return c.JSON(http.StatusOK, response)
}
