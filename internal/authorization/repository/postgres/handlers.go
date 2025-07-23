package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"
	"github.com/kkvaleriy/istokAuthorization/internal/authorization/repository/postgres/querys"
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

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, querys.AddUser, args)
	if err != nil {
		if pgxErr, ok := err.(pgx.PgError); ok && pgxErr.Code == "23505" {
			return fmt.Errorf("Not uniq user: %w", err)
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
