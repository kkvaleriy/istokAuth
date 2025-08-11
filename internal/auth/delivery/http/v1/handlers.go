package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	httperrors "github.com/kkvaleriy/istokAuth/internal/auth/delivery/http/v1/errors"
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
// @Failure 400 {object} httperrors.badRequestErrorResponse "Bad request"
// @Failure 409 {object} httperrors.validationDTOErrorResponse "A user already exists"
// @Failure 422 {object} httperrors.validationErrorResponse "Bad json in request"
// @Failure 500 {object} httperrors.internalServerErrorResponse "Internal server error"
// @Router /auth/signup [post]
func (h *handler) signUp(c echo.Context) error {

	request := &dtos.CreateUserRequest{}

	if err := c.Bind(request); err != nil {
		return &httperrors.BadRequestError{Err: err}
	}

	h.log.Debug("Request for signup", "host", c.Request().Host, "URL", c.Request().URL, "from", c.RealIP(), "request", request)

	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return httperrors.ErrValidation(err)
	}

	response, err := h.usecase.SignUp(c.Request().Context(), request)
	if err != nil {
		var validationError *dtos.ValidationError
		if errors.As(err, &validationError) {
			return &httperrors.ValidationDTOError{Err: validationError}
		}

		return err
	}

	h.log.Debug("Successfully signup", "host", c.Request().Host, "URL", c.Request().URL, "from", c.RealIP(), "request", request)
	h.log.Info("New user has been created", "user", request.Nickname, "from", c.RealIP())

	return c.JSON(http.StatusCreated, response)
}

// @Summary User authorization
// @Tags Authorization
// @Description User authorization
// @Accept json
// @Produce json
// @Param input body dtos.SignInRequest true "Account information for signin"
// @Success 200 {object} dtos.SignInResponse "Json with JWT, refresh token in coockie"
// @Failure 400 {object} httperrors.badRequestErrorResponse "Bad request"
// @Failure 401 {object} httperrors.authErrorResponse "Invalid credentials"
// @Failure 422 {object} httperrors.validationErrorResponse "Bad json in request"
// @Failure 500 {object} httperrors.internalServerErrorResponse "Internal server error"
// @Router /auth/signin [post]
func (h *handler) signIn(c echo.Context) error {
	request := &dtos.SignInRequest{}

	if err := c.Bind(request); err != nil {
		return &httperrors.BadRequestError{Err: err}
	}

	h.log.Debug("Request for signin", "host", c.Request().Host, "URL", c.Request().URL, "from", c.RealIP(), "request", request)

	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		return httperrors.ErrValidation(err)
	}

	response, err := h.usecase.SignIn(c.Request().Context(), request)
	if err != nil {
		var validationError *dtos.SignInError
		if errors.As(err, &validationError) {
			return &httperrors.AuthError{Err: validationError}
		}
		return err
	}

	coockie := &http.Cookie{
		Name:     "refreshToken",
		Value:    response.RToken,
		Expires:  response.ExpiresRToken,
		Path:     "/api/v1/auth",
		HttpOnly: true,
	}

	h.log.Debug("Successfully signin", "host", c.Request().Host, "URL", c.Request().URL, "from", c.RealIP(), "request", request)
	h.log.Info("The user has been successfully authenticated", "email", request.Email, "phone", request.Phone, "from", c.RealIP())

	c.SetCookie(coockie)
	return c.JSON(http.StatusOK, response)
}

// @Summary Refresh tokens
// @Tags Authorization
// @Description Get new refresh and access tokens by refresh token cookie
// @Produce json
// @Success 200 {object} dtos.SignInResponse "Json with JWT, refresh token in coockie"
// @Failure 400 {object} httperrors.badRequestErrorResponse "Bad request"
// @Failure 401 {object} httperrors.authErrorResponse "Invalid user"
// @Failure 500 {object} httperrors.internalServerErrorResponse "Internal server error"
// @Router /auth/refresh [get]
func (h *handler) Refresh(c echo.Context) error {
	h.log.Debug("Request for refresh tokens", "host", c.Request().Host, "URL", c.Request().URL, "from", c.RealIP())

	refreshTokenCookie, err := c.Request().Cookie("refreshToken")
	if err != nil {
		return &httperrors.BadRequestError{Err: errors.New("cookie refreshToken not set")}
	}

	if refreshTokenCookie.Expires.After(time.Now()) {
		return &httperrors.BadRequestError{Err: errors.New("cookie refreshToken is expiret")}
	}

	request, err := dtos.RequestByUUID(refreshTokenCookie.Value)
	if err != nil {
		return &httperrors.BadRequestError{errors.New("bad value of the refreshToken cookie")}
	}

	response, err := h.usecase.Refresh(c.Request().Context(), request)
	if err != nil {
		var validationError *dtos.SignInError
		if errors.As(err, &validationError) {
			return &httperrors.AuthError{Err: validationError}
		}
		return err
	}

	coockie := &http.Cookie{
		Name:     "refreshToken",
		Value:    response.RToken,
		Expires:  response.ExpiresRToken,
		Path:     "/api/v1/auth",
		HttpOnly: true,
	}

	h.log.Debug("Tokens have been successfully updated", "host", c.Request().Host, "URL", c.Request().URL, "from", c.RealIP())
	h.log.Info("The user has successfully updated the tokens", "from", c.RealIP())

	c.SetCookie(coockie)
	return c.JSON(http.StatusOK, response)
}

// @Summary User update password
// @Tags User
// @Description User update password
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body dtos.UpdateUserPasswordRequest true "New password"
// @Success 200
// @Failure 400 {object} httperrors.badRequestErrorResponse "Bad request"
// @Failure 401 {object} httperrors.authErrorResponse "Invalid token"
// @Failure 422 {object} httperrors.validationErrorResponse "Bad json in request"
// @Failure 500 {object} httperrors.internalServerErrorResponse "Internal server error"
// @Router /user/update-password [put]
func (h *handler) UpdateUserPassword(c echo.Context) error {
	user := c.Get("user").(jwt.MapClaims)

	h.log.Debug("Request for update password", "from", c.RealIP())

	sub, err := user.GetSubject()
	if err != nil {
		return &httperrors.BadRequestError{}
	}

	userUUID, err := dtos.RequestByUUID(sub)
	if err != nil {
		return &httperrors.BadRequestError{}
	}

	request := &dtos.UpdateUserPasswordRequest{}

	if err := c.Bind(request); err != nil {
		return &httperrors.BadRequestError{Err: err}
	}

	validate := validator.New()

	err = validate.Struct(request)
	if err != nil {
		return httperrors.ErrValidation(err)
	}

	request.ID = userUUID

	err = h.usecase.UpdateUserPassword(c.Request().Context(), request)
	if err != nil {
		var validationError *dtos.ValidationError
		if errors.As(err, &validationError) {
			return &httperrors.ValidationDTOError{Err: validationError}
		}

		return err
	}

	h.log.Info("The user has successfully updated password", "from", c.RealIP())

	return c.NoContent(http.StatusOK)
}
