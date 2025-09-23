package memory

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type UserStore struct {
	mu sync.RWMutex

	idSeq  atomic.Int64
	users  map[int64]model.User
	emails map[string]int64

	now func() time.Time
}

func NewUserStore() *UserStore {
	const size = 256

	store := &UserStore{
		users:  make(map[int64]model.User, size),
		emails: make(map[string]int64, size),
		now:    func() time.Time { return time.Now().UTC() },
	}

	store.idSeq.Store(0)

	return store
}

func (s *UserStore) Create(ctx context.Context, user model.User) (model.User, error) {
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
	}

	if err := s.isEmailExists(user.Email, s.mu.RLocker()); err != nil {
		return model.User{}, err
	}

	id := s.nextID()

	user.ID = id
	user.CreateTime = s.now()

	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.isEmailExists(user.Email, nil); err != nil {
		return model.User{}, err
	}

	s.users[id] = user
	s.emails[user.Email] = id

	return user, nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (model.User, error) {
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return model.User{}, model.ErrNotFound
	}

	return user, nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (model.User, error) {
	select {
	case <-ctx.Done():
		return model.User{}, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.emails[email]
	if !ok {
		return model.User{}, model.ErrNotFound
	}

	return s.users[id], nil
}

func (s *UserStore) nextID() int64 {
	return s.idSeq.Add(1)
}

func (s *UserStore) isEmailExists(email string, locker sync.Locker) error {
	if locker != nil {
		locker.Lock()
		defer locker.Unlock()
	}

	if _, ok := s.emails[email]; ok {
		return model.ErrAlreadyExists
	}

	return nil
}
