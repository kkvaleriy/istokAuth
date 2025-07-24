package postgres

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kkvaleriy/istokAuthorization/internal/authorization/dtos"
	user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"
)

func createUserArgs(u *user.User) pgx.NamedArgs {
	return pgx.NamedArgs{
		"name":      u.Name,
		"lastname":  u.Lastname,
		"nickname":  u.Nickname,
		"email":     u.Email,
		"userType":  u.UserType,
		"isActive":  u.IsActive,
		"phone":     u.Phone,
		"UUID":      u.UUID,
		"passHash":  u.PassHash[:],
		"createdAt": u.CreatedAt,
	}
}

func ErrorValidation(constraint string, args pgx.NamedArgs) *dtos.ValidationError {
	var field, value string

	switch strings.ToLower(constraint) {
	case "uniq_email":
		field = "email"
	case "uniq_nickname":
		field = "nickname"
	case "uniq_phone":
		field = "phone"
	}

	if len(field) < 1 {
		return &dtos.ValidationError{
			Message: "unknown error",
		}
	}

	switch args[field].(type) {
	case int:
		value = strconv.Itoa(args[field].(int))
	case string:
		value = args[field].(string)
	}

	return &dtos.ValidationError{
		Message: fmt.Sprintf("a user with the %s %s already exists", field, value),
		Field:   field,
		Value:   value,
	}
}
