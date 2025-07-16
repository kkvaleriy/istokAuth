package v1

import (
	"net/http"

	"github.com/kkvaleriy/istokAuthorization/internal/authorization/dtos"
	"github.com/labstack/echo/v4"
)

type handler struct {
	usecase Usecase
}

