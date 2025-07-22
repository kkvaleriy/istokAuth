package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	user "github.com/kkvaleriy/istokAuthorization/internal/authorization/entities"
	"github.com/kkvaleriy/istokAuthorization/internal/authorization/repository/postgres/querys"
)

type repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) AddUser(ctx context.Context, u *user.User) error {
	args := createUserArgs(u)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// TODO: check uniq 23505
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
