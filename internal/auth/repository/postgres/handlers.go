package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	user "github.com/kkvaleriy/istokAuth/internal/auth/entities"
	"github.com/kkvaleriy/istokAuth/internal/auth/repository/postgres/querys"
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

	return r.insertInDB(ctx, querys.AddUser, args)
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
