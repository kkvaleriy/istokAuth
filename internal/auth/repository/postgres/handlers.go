package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kkvaleriy/istokAuth/internal/auth/dtos"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
	"github.com/kkvaleriy/istokAuth/internal/auth/repository/postgres/queries"
)

type logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type repository struct {
	db  *pgxpool.Pool
	log logger
}

func New(db *pgxpool.Pool, log logger) *repository {
	return &repository{db: db, log: log}
}

func (r *repository) AddUser(ctx context.Context, u *user.User) error {
	args := createUserArgs(u)
	r.log.Debug("the args for query have been created", "args", args)

	return r.insertInDB(ctx, queries.AddUser, args)
}

func (r *repository) CheckUserByCredentials(ctx context.Context, u *user.User) (*user.User, error) {
	args := checkUserByCredentialsArgs(u)

	err := r.db.QueryRow(ctx, queries.CheckUserByCredentials, args).Scan(&u.UUID, &u.Nickname, &u.UserType, &u.IsActive)
	if err != nil {
		return nil, signInError(err)
	}

	return u, nil
}

func (r *repository) AddToken(ctx context.Context, t *user.RToken) error {
	args := addTokenArgs(t)
	r.log.Debug("the args for query have been created", "args", args)

	return r.insertInDB(ctx, queries.AddRToken, args)

}

func (r *repository) insertInDB(ctx context.Context, query string, args pgx.NamedArgs) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	r.log.Debug("the transaction has started")

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		r.log.Debug("the transaction was failed, rollback")
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return errorValidation(pgErr.ConstraintName, args)
		}
		return err
	}

	r.log.Debug("the request was successful", "query", query, "args", args)

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	r.log.Debug("the transaction has finished successfuly")

	return nil
}

func (r *repository) RefreshToken(ctx context.Context, u *user.User, t *user.RToken) (*user.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	r.log.Debug("the transaction has started")

	err = tx.QueryRow(ctx, queries.DelRToken, t.UUID).Scan(&u.UUID, &u.Nickname)
	if err != nil {
		return nil, err
	}
	if len(u.UUID) < 1 {
		return nil, &dtos.SignInError{Message: "Invalid refresh token"}
	}

	err = tx.QueryRow(ctx, queries.UserType, u.UUID).Scan(&u.UserType)
	if err != nil {
		return nil, err
	}
	if len(u.UserType) < 1 {
		return nil, &dtos.SignInError{Message: "Invalid user"}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return u, nil
}
