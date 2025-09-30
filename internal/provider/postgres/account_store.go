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
	//go:embed queries/create_account.sql
	createAccountQuery string

	//go:embed queries/get_account_by_id.sql
	getAccountByIDQuery string

	//go:embed queries/get_accounts_by_user_id.sql
	getAccountsByUserIDQuery string
)

type AccountStore struct {
	db  *pgxpool.Pool
	now func() time.Time
}

func NewAccountStore(db *pgxpool.Pool) *AccountStore {
	return &AccountStore{
		db:  db,
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (s *AccountStore) Create(ctx context.Context, account model.Account) (model.Account, error) {
	row := s.db.QueryRow(ctx, createAccountQuery, account.UserID, account.Title, account.IsDefault, s.now())

	err := row.Scan(&account.ID, &account.CreateTime)
	if err != nil {
		return model.Account{}, fmt.Errorf("insert account: %w", err)
	}

	return account, nil
}

func (s *AccountStore) GetByID(ctx context.Context, id int64) (model.Account, error) {
	row := s.db.QueryRow(ctx, getAccountByIDQuery, id)

	account := model.Account{}

	err := row.Scan(&account.ID, &account.UserID, &account.Title, &account.IsDefault, &account.CreateTime)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Account{}, fmt.Errorf("account not found: %w", model.ErrNotFound)
	}
	if err != nil {
		return model.Account{}, fmt.Errorf("select account by id: %w", err)
	}

	return account, nil
}

func (s *AccountStore) GetAccountsByUserID(ctx context.Context, userID int64) ([]model.Account, error) {
	rows, err := s.db.Query(ctx, getAccountsByUserIDQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("select accounts by user id: %w", err)
	}

	accounts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Account, error) {
		var account model.Account

		err = row.Scan(&account.ID, &account.UserID, &account.Title, &account.IsDefault, &account.CreateTime)
		if err != nil {
			return model.Account{}, err
		}

		return account, nil
	})
	if err != nil {
		return nil, fmt.Errorf("select accounts by user id: %w", err)
	}

	return accounts, nil
}
