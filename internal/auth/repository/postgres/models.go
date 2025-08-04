package postgres

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
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

func checkUserByCredentialsArgs(u *user.User) pgx.NamedArgs {
	return pgx.NamedArgs{
		"email":    u.Email,
		"phone":    u.Phone,
		"passHash": u.PassHash,
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
