package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/meetmorrowsolonmars/education-pet-project/internal/domain/model"
)

type OperationStore struct {
	mu sync.RWMutex

	operations []model.Operation

	newID func() (uuid.UUID, error)
	now   func() time.Time
}

func NewOperationStore() *OperationStore {
	const size = 256

	return &OperationStore{
		operations: make([]model.Operation, 0, size),
		newID:      func() (uuid.UUID, error) { return uuid.NewV7() },
		now:        func() time.Time { return time.Now().UTC() },
	}
}

func (s *OperationStore) Create(ctx context.Context, operation model.Operation) (model.Operation, error) {
	select {
	case <-ctx.Done():
		return model.Operation{}, ctx.Err()
	default:
	}

	id, err := s.newID()
	if err != nil {
		return model.Operation{}, fmt.Errorf("generate id: %w", err)
	}

	operation.ID = id
	operation.CreateTime = s.now()

	s.mu.Lock()
	defer s.mu.Unlock()

	s.operations = append(s.operations, operation)

	return operation, nil
}

func (s *OperationStore) GetByAccountID(
	ctx context.Context,
	accountID int64,
	limit int64,
	offset int64,
) ([]model.Operation, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()

	operations := make([]model.Operation, 0)

	for _, operation := range s.operations {
		if operation.AccountID == accountID {
			operations = append(operations, operation)
		}
	}

	s.mu.RUnlock()

	sort.Slice(operations, func(i, j int) bool {
		return operations[i].CreateTime.After(operations[j].CreateTime)
	})

	if len(operations) == 0 {
		return []model.Operation{}, nil
	}

	if offset >= int64(len(operations)) {
		return []model.Operation{}, nil
	}

	end := offset + limit
	if end > int64(len(operations)) {
		end = int64(len(operations))
	}

	return operations[offset:end], nil
}
