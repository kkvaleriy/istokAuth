package postgres

import (
	"github.com/jackc/pgx/v5"
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
