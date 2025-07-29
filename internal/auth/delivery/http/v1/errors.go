package v1

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type HTTPError interface {
	error
	StatusCode() int
}

type BadRequestError struct {
	Err error
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("Bad json in request: %s", e.Err.Error())
}

func (e *BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}

// @Description JSON validate error
type validationErrorResponse struct {
	Error  string   `json:"error" example:"Bad json in request"`
	Fields []string `json:"fields" example:"email:required"`
}

type ValidationError struct {
	Fields []string
}

func ErrValidation(err error) *ValidationError {
	errValidation, ok := err.(validator.ValidationErrors)
	if !ok {
		return &ValidationError{}
	}
	fields := []string{}
	for _, e := range errValidation {
		fields = append(fields, fmt.Sprintf("%s:%s", e.Field(), e.Tag()))
	}
	return &ValidationError{Fields: fields}
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func (e *ValidationError) StatusCode() int {
	return http.StatusBadRequest
}

// @Description Uniqueness error
type validationDTOErrorResponse struct {
	Error string `json:"error" example:"A user with the nickname Johny1 already exists"`
}

type ValidationDTOError struct {
	Err error
}

func (e *ValidationDTOError) Error() string {
	return e.Err.Error()
}

func (e *ValidationDTOError) StatusCode() int {
	return http.StatusConflict
}

// @Description Internal server error
type internalServerErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}

func ErrorsHandler(log logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var status int
		var message interface{}

		switch e := err.(type) {
		case *BadRequestError:
			status = e.StatusCode()
			message = map[string]interface{}{
				"error": e.Error()}
		case *ValidationError:
			status = e.StatusCode()
			message = map[string]interface{}{
				"error":  e.Error(),
				"fields": e.Fields,
			}
		case *ValidationDTOError:
			status = e.StatusCode()
			message = map[string]interface{}{
				"error": e.Error()}
		case *echo.HTTPError:
			status = e.Code
			message = map[string]interface{}{
				"error": e.Error()}
		default:
			status = http.StatusInternalServerError
			message = map[string]interface{}{
				"error": "Internal server error"}
		}

		log.Error("request error", "error", err.Error(), "status", status)

		if err := c.JSON(status, message); err != nil {
			log.Error("failed to send error response", "error", err.Error())
		}
	}
}
