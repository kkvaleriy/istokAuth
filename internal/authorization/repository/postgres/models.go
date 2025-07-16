package postgres

import (
	"time"

	"github.com/google/uuid"
	user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"
)

type createUserModel struct {
	Name      string    `db:"name"`
	Lastname  string    `db:"lastname"`
	Nickname  string    `db:"nickname"`
	Email     string    `db:"email"`
	UserType  string    `db:"userType"`
	IsActive  bool      `db:"isActive"`
	Phone     int       `db:"phone"`
	UUID      uuid.UUID `db:"UUID"`
	PassHash  [32]byte  `db:"passHash"`
	CreatedAt time.Time `db:"createdAt"`
}

func newCreateUserModel(u *user.User) *createUserModel {
	return &createUserModel{
		Name:      u.Name,
		Lastname:  u.Lastname,
		Nickname:  u.Nickname,
		Email:     u.Nickname,
		UserType:  u.UserType,
		IsActive:  u.IsActive,
		Phone:     u.Phone,
		UUID:      u.UUID,
		PassHash:  u.PassHash,
		CreatedAt: u.CreatedAt,
	}
}
