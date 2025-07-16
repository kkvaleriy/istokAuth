package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"
	"github.com/kkvaleriy/istokAuthorization/internal/authorization/repository/postgres/querys"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) AddUser(ctx context.Context, u *user.User) error {

	userModel := newCreateUserModel(u)

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(querys.AddUser, userModel)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
