package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

var (
	//go:embed queries/create_user.sql
	createUserQuery string

	//go:embed queries/get_user_by_id.sql
	getUserByIDQuery string

	//go:embed queries/get_user_by_email.sql
	getUserByEmailQuery string
)

type UserStore struct {
	db  *pgxpool.Pool
	now func() time.Time
}

func NewUserStore(db *pgxpool.Pool) *UserStore {
	return &UserStore{
		db:  db,
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (s *UserStore) Create(ctx context.Context, user model.User) (model.User, error) {
	err := pgx.BeginFunc(ctx, s.db, func(tx pgx.Tx) error {
		now := s.now()
		row := tx.QueryRow(ctx, createUserQuery, user.Email, user.Password, user.FullName, now)

		err := row.Scan(&user.ID, &user.CreateTime)
		if err != nil {
			return fmt.Errorf("insert user: %w", err)
		}

		_, err = tx.Exec(ctx, createAccountQuery, user.ID, "default", true, now)
		if err != nil {
			return fmt.Errorf("create default account: %w", err)
		}

		return nil
	})
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (model.User, error) {
	row := s.db.QueryRow(ctx, getUserByIDQuery, id)

	user := model.User{}

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.CreateTime)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf("user not found: %w", model.ErrNotFound)
	}
	if err != nil {
		return model.User{}, fmt.Errorf("select user by id: %w", err)
	}

	return user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (model.User, error) {
	row := s.db.QueryRow(ctx, getUserByEmailQuery, email)

	user := model.User{}

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FullName, &user.CreateTime)
	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, fmt.Errorf("user not found: %w", model.ErrNotFound)
	}
	if err != nil {
		return model.User{}, fmt.Errorf("select user by email: %w", err)
	}

	return user, nil
}
