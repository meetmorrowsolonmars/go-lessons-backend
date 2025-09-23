package user

import (
	"context"
	"fmt"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type Service struct {
	store        Store
	accountStore AccountStore
}

func NewService(
	store Store,
	accountStore AccountStore,
) *Service {
	return &Service{
		store:        store,
		accountStore: accountStore,
	}
}

func (s *Service) Create(ctx context.Context, user model.User) (model.User, error) {
	var err error

	err = user.HashPassword()
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err = s.store.Create(ctx, user)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	_, err = s.accountStore.CreateDefault(ctx, user.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("create default account: %w", err)
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (model.User, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) GetByEmail(ctx context.Context, email string) (model.User, error) {
	return s.store.GetByEmail(ctx, email)
}
