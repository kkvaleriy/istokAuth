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

type JSONParseError struct {
	Err error
}

func (e *JSONParseError) Error() string {
	return fmt.Sprintf("Bad json in request: %s", e.Err.Error())
}

func (e *JSONParseError) StatusCode() int {
	return http.StatusBadRequest
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

type ValidationErrorDTO struct {
	Err error
}

func (e *ValidationErrorDTO) Error() string {
	return e.Err.Error()
}

func (e *ValidationErrorDTO) StatusCode() int {
	return http.StatusConflict
}

func ErrorsHandler(log logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		var status int
		var message interface{}

		switch e := err.(type) {
		case *JSONParseError:
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
