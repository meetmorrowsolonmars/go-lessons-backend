package memory

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type AccountStore struct {
	mu sync.RWMutex

	idSeq    atomic.Int64
	accounts map[int64]model.Account

	now func() time.Time
}

func NewAccountStore() *AccountStore {
	const size = 256

	store := &AccountStore{
		accounts: make(map[int64]model.Account, size),
		now:      func() time.Time { return time.Now().UTC() },
	}

	store.idSeq.Store(0)

	return store
}

func (s *AccountStore) Create(ctx context.Context, account model.Account) (model.Account, error) {
	select {
	case <-ctx.Done():
		return model.Account{}, ctx.Err()
	default:
	}

	id := s.nextID()

	account.ID = id
	account.CreateTime = s.now()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.accounts[id] = account

	return account, nil
}

func (s *AccountStore) CreateDefault(ctx context.Context, userID int64) (model.Account, error) {
	account := model.Account{
		UserID:    userID,
		Title:     "default",
		IsDefault: true,
	}

	return s.Create(ctx, account)
}

func (s *AccountStore) GetByID(ctx context.Context, id int64) (model.Account, error) {
	select {
	case <-ctx.Done():
		return model.Account{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	account, ok := s.accounts[id]
	if !ok {
		return model.Account{}, model.ErrNotFound
	}

	return account, nil
}

func (s *AccountStore) GetAccountsByUserID(ctx context.Context, userID int64) ([]model.Account, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	accounts := make([]model.Account, 0, 1)

	for _, account := range s.accounts {
		if account.UserID == userID {
			accounts = append(accounts, account)
		}
	}

	return accounts, nil
}

func (s *AccountStore) nextID() int64 {
	return s.idSeq.Add(1)
}
